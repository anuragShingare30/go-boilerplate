package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/anuragShingare30/go-boilerplate/internal/config"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrzerolog"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

/**
@dev NewRelic: Tracks response times, throughput, and error rates,
Identifies slow database queries and API calls, Shows which endpoints are consuming most resources

@dev logger Service: Collects all application logs in one central place.
And, Makes debugging production issues much easier

@dev Integration flow: (imp!!!!!)

NewLoggerService (initializes NewRelic) 
    → NewLoggerWithService (creates logger with NewRelic integration)
        → Application code uses the logger
            → Logs automatically forwarded to NewRelic in production
*/

// here struct element is in small case - internal use only
type LoggerService struct {
	nrApp *newrelic.Application
}

// NewLoggerService: initializes and returns a new LoggerService instance
// we use newrelic package to init the NewRelic service/application
func NewLoggerService(cfg *config.ObservabilityConfig) *LoggerService {
	service := &LoggerService{
		nrApp: nil,
	}

	if cfg.NewRelic.LicenseKey == "" {
		return service
	}

	var configOptions []newrelic.ConfigOption
	configOptions = append(configOptions,
		newrelic.ConfigAppName(cfg.ServiceName),
		newrelic.ConfigLicense(cfg.NewRelic.LicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(cfg.NewRelic.AppLogForwardingEnabled),
		newrelic.ConfigDistributedTracerEnabled(cfg.NewRelic.DistributedTracingEnabled),
	)

	// Add debug logging only if explicitly enabled in observability config
	if cfg.NewRelic.DebugLogging {
		configOptions = append(configOptions,
			newrelic.ConfigDebugLogger(os.Stdout),
		)
	}

	// initialize the new logger service
	app, err := newrelic.NewApplication(configOptions...)
	if err != nil {
		fmt.Println("failed to initialized logger service")
		return service
	}

	service.nrApp = app
	fmt.Println("Successfully initialized logger service")
	return service
}

// Shutdown: Gracefully shuts down New Relic
func (ls *LoggerService) Shutdown() {
	if ls.nrApp != nil {
		ls.nrApp.Shutdown(10 * time.Second)
	}
}

// GetApplication: returns the New Relic application instance
func (ls *LoggerService) GetApplication() *newrelic.Application {
	return ls.nrApp
}

// NewLoggerWithService creates logger with NewRelic integration
func NewLoggerWithService(cfg *config.ObservabilityConfig, loggerService *LoggerService) zerolog.Logger {
	var logLevel zerolog.Level
	level := cfg.GetLogLevel()

	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	// Don't set global level - let each logger have its own level
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var writer io.Writer

	// Setup base writer
	// If LoggerService has an active NewRelic app, wraps the writer with zerologWriter.New() to automatically forward logs to NewRelic
	if cfg.IsProduction() && cfg.Logging.Format == "json" {
		// In production, write to stdout
		writer = os.Stdout

		// Wrap with New Relic zerologWriter for log forwarding in production
		// ....
	} else {
		// Uses ConsoleWriter for human-readable, colored output
		// No NewRelic integration (logs stay local)
		// Development mode - use console writer
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
		writer = consoleWriter
	}

	// Logger creation
	logger := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Str("environment", cfg.Environment).
		Logger()

	// Include stack traces for errors in development
	if !cfg.IsProduction() {
		logger = logger.With().Stack().Logger()
	}

	// Add New Relic hook for log forwarding in production
	if cfg.IsProduction() && loggerService != nil && loggerService.nrApp != nil {
		nrHook := nrzerolog.NewRelicHook{
			App: loggerService.nrApp,
		}
		logger = logger.Hook(nrHook)
	}

	return logger
}


// WithTraceContext: adds New Relic transaction context to logger
// newrelic.Transaction: represents a single web request or background task being monitored by NewRelic. 
// It's typically created at the start of an HTTP handler using the NewRelic middleware.
// Request duration, Response status codes, Database query times, External API calls, Errors and panics, Custom events/metrics
// kind of trace which have a starting point and end point, all the interaction and components it touches during request lifecylce are included in single transaction. If something goes wrong, we can take a particular tnx and explore.
func WithTraceContext(logger zerolog.Logger, txn *newrelic.Transaction) zerolog.Logger {
	if txn == nil {
		return logger
	}

	metadata := txn.GetTraceMetadata()

	return logger.With().
		Str("trace.id", metadata.TraceID).
		Str("span.id", metadata.SpanID).Logger()
}