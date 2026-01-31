package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/anuragShingare30/go-boilerplate/internal/config"
	loggerConfig "github.com/anuragShingare30/go-boilerplate/internal/logger"
	pgxzero "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/newrelic/go-agent/v3/integrations/nrpgx5"
	"github.com/rs/zerolog"
)

// @dev logic to connect the db
// @dev Database Pooling: Opening certain number of connections in-hand, application will be more efficient in performance


type Database struct {
	Pool *pgxpool.Pool // to store pool
	log *zerolog.Logger // to log db related info
}

type multiTracer struct{
	tracers []any
}

// DatabaseTimeout is the timeout duration for database operations in seconds.
const DatabasePingTimeout = 10

// TraceQueryStart implements pgx tracer interface
func (mt *multiTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryStart(context.Context, *pgx.Conn, pgx.TraceQueryStartData) context.Context
		}); ok {
			ctx = t.TraceQueryStart(ctx, conn, data)
		}
	}
	return ctx
}

// TraceQueryEnd implements pgx tracer interface
func (mt *multiTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryEnd(context.Context, *pgx.Conn, pgx.TraceQueryEndData)
		}); ok {
			t.TraceQueryEnd(ctx, conn, data)
		}
	}
}


func New(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerConfig.LoggerService) (*Database, error){
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

	// URL-encode the password
	encodedPassword := url.QueryEscape(cfg.Database.Password)
	// format the string and return new string
	// postgres://user:password@host:port/dbname?sslmode=disable
	dns := fmt.Sprintf("postgres://%s:%s@%s/dbname?sslmode=%s",
		cfg.Database.User,
		encodedPassword,
		hostPort,
		cfg.Database.SSLMode,
	)

	// Converts DSN into pgxpool.Config (connection pool settings)
	pgxPoolConfig, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	
	// Add New Relic PostgreSQL instrumentation
	if loggerService != nil && loggerService.GetApplication() != nil {
		pgxPoolConfig.ConnConfig.Tracer = nrpgx5.NewTracer()
	}

	// Development: you want to see SQL queries in your console
	// Production:  you only want them in New Relic and not in console
	if cfg.Primary.Env == "local" {
		globalLevel := logger.GetLevel()
		pgxLogger := loggerConfig.NewPgxLogger(globalLevel)
		// Chain tracers - New Relic first, then local logging
		if pgxPoolConfig.ConnConfig.Tracer != nil {
			// Creates a local tracer
			localTracer := &tracelog.TraceLog{
				Logger:   pgxzero.NewLogger(pgxLogger),
				LogLevel: tracelog.LogLevel(loggerConfig.GetPgxTraceLogLevel(globalLevel)),
			}
			// multitracer: Newrelic + console logging
			pgxPoolConfig.ConnConfig.Tracer = &multiTracer{
				tracers: []any{pgxPoolConfig.ConnConfig.Tracer, localTracer},
			}
		} else {
			pgxPoolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
				Logger:   pgxzero.NewLogger(pgxLogger),
				LogLevel: tracelog.LogLevel(loggerConfig.GetPgxTraceLogLevel(globalLevel)),
			}
		}
	}
	

	// Establishes actual database connections
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil{
		return nil, fmt.Errorf("failed to establish database connection %w", err)
	}

	database := &Database{
		Pool: pool,
		log: logger,
	}

	// Pings database with 10-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), DatabasePingTimeout*time.Second)
	defer cancel()
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info().Msg("connected to database!!!")


	return database, nil
}


// Close: gracefully closes the database connection pool
func (db *Database) Close() error {
	db.log.Info().Msg("closing database connection pool!!!")
	db.Pool.Close()
	return nil
}