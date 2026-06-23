package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App 			App 			    `mapstructure:"app"`
	Env            	string              `mapstructure:"env"`
	HTTP           	HTTPConfig          `mapstructure:"http"`
	Postgres       	PostgresConfig      `mapstructure:"postgres"`
	Redis          	RedisConfig         `mapstructure:"redis"`
	Cache          	CacheConfig         `mapstructure:"cache"`
	GRPC           	GRPCConfig          `mapstructure:"grpc"`
	Logging        	Logging             `mapstructure:"logging"`
	Observability  	ObservabilityConfig `mapstructure:"observability"`
}

type App struct{
	ShutdownTimeout string `mapstructure:"shutdown_timeout"`
}

type Logging struct {
	File   string `mapstructure:"file"`
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type ObservabilityConfig struct {
	OTLPEndpoint string `mapstructure:"otlp_endpoint"`
	ServiceName  string `mapstructure:"service_name"`
}

type HTTPConfig struct {
	Host string `mapstructure:"host"`
	Port int 	`mapstructure:"port"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	MaxConns int 	`mapstructure:"max_conns"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type CacheConfig struct {
	QuestTTL time.Duration `mapstructure:"quest_ttl"`
}

type GRPCConfig struct {
	NotificationAddr string `mapstructure:"notification_addr"`
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.DBName, p.SSLMode,
	)
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	v.SetEnvPrefix("QUESTHUB")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) ShutdownTimeoutDuration() (time.Duration, error) {
	return time.ParseDuration(c.App.ShutdownTimeout)
}

func (c *Config) CacheQuestTTLDuration() (time.Duration, error) {
	return time.ParseDuration(c.Cache.QuestTTL.String())
}
