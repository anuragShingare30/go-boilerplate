package config

// importing packages
import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

// @dev to load all env variables in struct when server starts
// @dev this loads all env variable into struct


type Config struct {
	Primary  Primary        `koanf:"primary" validation:"required"`
	Server   ServerConfig   `koanf:"server" validation:"required"`
	Redis    RedisConfig    `koanf:"redis" validation:"required"`
	Database DatabaseConfig `koanf:"database" validation:"required"`
	Auth     AuthConfig     `koanf:"auth" validation:"required"`
	Observability *ObservabilityConfig `koanf:"observability" validation:"required"`
}

type Primary struct {
	Env string `koanf:"env" validation:"required"`
}

type ServerConfig struct {
	Port               string   `koanf:"port" validation:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validation:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validation:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validation:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validation:"required"`
}

type RedisConfig struct {
	Address string `koanf:"address" validation:"required"`
}

type DatabaseConfig struct {
	Host            string `koanf:"host" validation:"required"`
	Port            int    `koanf:"port" validation:"required"`
	User            string `koanf:"user" validation:"required"`
	Password        string `koanf:"password" validation:"required"`
	Name            string `koanf:"name" validation:"required"`
	SSLMode         string `koanf:"ssl_mode" validation:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns" validation:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validation:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validation:"required"`
	ConnMaxIdletime int    `koanf:"conn_max_idletime" validation:"required"`
}

type AuthConfig struct {
	SecretKey string `koanf:"secret_key" validation:"required"`
}

// LoadConfig loads the configuration from environment variables using koanf
func LoadConfig() (mainConfig *Config, err error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	// loading env variables using koanf
	k := koanf.New(".")

	err = k.Load(env.Provider("BOILERPLATE_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "BOILERPLATE_"))
	}), nil)
	// err != nil -> checks if error exists
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig = &Config{}

	err = k.Unmarshal("", mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal mainconfig")
	}

	validate := validator.New()

	err = validate.Struct(mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not validate the struct")
	}

	// set default observability config if not provided
	// in config struct we set Observability as pointer type to check whether it is nil or not
	if mainConfig.Observability == nil {
		mainConfig.Observability = DefaultObservabilityConfig()
	}

	// fill some of the fields
	mainConfig.Observability.ServiceName = "go-boilerplate"
	mainConfig.Observability.Environment = mainConfig.Primary.Env

	// automatic pointer dereferencing for method calls
	err = mainConfig.Observability.Validate()
	if err != nil {
		logger.Fatal().Err(err).Msg("invalid observability config")
	}

	return
}
