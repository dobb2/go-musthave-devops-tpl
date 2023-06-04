package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/caarlos0/env/v7"
)

type Config struct {
	Address         string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreFile       string        `env:"STORE_FILE" envDefault:"tmp/devops-metrics-db.json"`
	Key             string        `env:"KEY" envDefault:""`
	DatabaseDSN     string        `env:"DATABASE_DSN" envDefault:""`
	ReportInterval  time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval    time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	Restore         bool          `env:"RESTORE" envDefault:"true"`
	RateLimit       int           `env:"RATE_LIMIT" envDefault:"0"`
	MetricMaxAmount int           `env:"METRIC_MAX_AMOUNT" envDefault:"37"`
}

type AgentConfig struct {
	Address         string
	Key             string
	ReportInterval  time.Duration
	PollInterval    time.Duration
	RateLimit       int
	MetricMaxAmount int
}

type ServerConfig struct {
	Address       string
	Key           string
	StoreFile     string
	DatabaseDSN   string
	Restore       bool
	StoreInterval time.Duration
}

func CreateAgentConfig(logger zerolog.Logger) AgentConfig {
	var envcfg Config
	err := env.Parse(&envcfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed parsed env config")
	}

	var cfg AgentConfig

	flag.StringVar(&cfg.Address, "a", envcfg.Address, "a string")
	flag.DurationVar(&cfg.ReportInterval, "r", envcfg.ReportInterval, "a duration")
	flag.DurationVar(&cfg.PollInterval, "p", envcfg.PollInterval, "a duration")
	flag.StringVar(&cfg.Key, "k", envcfg.Key, "a string")
	flag.IntVar(&cfg.RateLimit, "l", envcfg.RateLimit, "a int")
	flag.IntVar(&cfg.MetricMaxAmount, "m", envcfg.MetricMaxAmount, "a int")

	flag.Parse()

	envStrAddres, boolAddres := os.LookupEnv("ADDRESS")
	if boolAddres {
		cfg.Address = envStrAddres
	}

	envStrKey, boolKey := os.LookupEnv("KEY")
	if boolKey {
		cfg.Key = envStrKey
	}

	envStrPoll, boolPoll := os.LookupEnv("POLL_INTERVAL")
	if boolPoll {
		envTimePoll, err := time.ParseDuration(envStrPoll)
		if err != nil {
			logger.Warn().Err(err).Msg("invalid poll interval time in env export")
		} else {
			cfg.PollInterval = envTimePoll
		}
	}

	envStrReport, boolReport := os.LookupEnv("REPORT_INTERVAL")
	if boolReport {
		envTimeReport, err := time.ParseDuration(envStrReport)
		if err != nil {
			logger.Warn().Err(err).Msg("invalid report interval time in env export")
		} else {
			cfg.PollInterval = envTimeReport
		}
	}

	envStrRateLimit, boolRateLimit := os.LookupEnv("RATE_LIMIT")
	if boolRateLimit {
		envRateLimit, err := strconv.Atoi(envStrRateLimit)
		if err != nil {
			logger.Warn().Err(err).Msg("invalid rate limit int in env export")
		} else {
			cfg.RateLimit = envRateLimit
		}
	}

	envStrMetricMaxAmount, boolMetricMaxAmount := os.LookupEnv("METRIC_MAX_AMOUNT")
	if boolMetricMaxAmount {
		envMetricMaxAmount, err := strconv.Atoi(envStrMetricMaxAmount)
		if err != nil {
			logger.Warn().Err(err).Msg("invalid metrics max amount int in env export")
		} else {
			cfg.MetricMaxAmount = envMetricMaxAmount
		}
	}

	return cfg
}

func CreateServerConfig(logger zerolog.Logger) ServerConfig {
	var envcfg Config
	err := env.Parse(&envcfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed parsed env config")
	}

	var cfg ServerConfig

	flag.StringVar(&cfg.Address, "a", envcfg.Address, "a string")
	flag.StringVar(&cfg.StoreFile, "f", envcfg.StoreFile, "file store a string")
	flag.StringVar(&cfg.DatabaseDSN, "d", envcfg.DatabaseDSN, "DSN a string")
	flag.StringVar(&cfg.Key, "k", envcfg.Key, "a string")
	flag.BoolVar(&cfg.Restore, "r", envcfg.Restore, "a bool")
	flag.DurationVar(&cfg.StoreInterval, "i", envcfg.StoreInterval, "a duration")
	flag.Parse()

	envStrAddres, boolAddres := os.LookupEnv("ADDRESS")
	if boolAddres {
		cfg.Address = envStrAddres
	}

	envStrKey, boolKey := os.LookupEnv("KEY")
	if boolKey {
		cfg.Key = envStrKey
	}

	envStrDSN, boolDSN := os.LookupEnv("DATABASE_DSN")
	if boolDSN {
		cfg.DatabaseDSN = envStrDSN
	}

	envStrFile, boolFile := os.LookupEnv("STORE_FILE")
	if boolFile {
		cfg.StoreFile = envStrFile
	}

	envStrRestore, boolRestore := os.LookupEnv("RESTORE")
	if boolRestore {
		envBoolRestore, err := strconv.ParseBool(envStrRestore)
		if err != nil {
			logger.Warn().Err(err).Msg("invalid restore bool in env export")
		} else {
			cfg.Restore = envBoolRestore
		}
	}

	envStrStore, boolStore := os.LookupEnv("STORE_INTERVAL")
	if boolStore {
		envTimeStore, err := time.ParseDuration(envStrStore)
		if err != nil {
			logger.Warn().Err(err).Msg("invalid store interval time in env export")
		} else {
			cfg.StoreInterval = envTimeStore
		}
	}

	return cfg
}
