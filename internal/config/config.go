package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v7"
)

type Config struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreFile      string        `env:"STORE_FILE" envDefault:"tmp/devops-metrics-db.json"`
	Key            string        `env:"KEY" envDefault:""`
	DatabaseDSN    string        `env:"DATABASE_DSN" envDefault:""`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	Restore        bool          `env:"RESTORE" envDefault:"true"`
}

type AgentConfig struct {
	Address        string
	Key            string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

type ServerConfig struct {
	Address       string
	Key           string
	StoreFile     string
	DatabaseDSN   string
	Restore       bool
	StoreInterval time.Duration
}

func CreateAgentConfig() AgentConfig {
	var envcfg Config
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}

	var cfg AgentConfig

	flag.StringVar(&cfg.Address, "a", envcfg.Address, "a string")
	flag.DurationVar(&cfg.ReportInterval, "r", envcfg.ReportInterval, "a duration")
	flag.DurationVar(&cfg.PollInterval, "p", envcfg.PollInterval, "a duration")
	flag.StringVar(&cfg.Key, "k", envcfg.Key, "a string")

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
			log.Println("invalid poll interval time in env export")
		} else {
			cfg.PollInterval = envTimePoll
		}
	}

	envStrReport, boolReport := os.LookupEnv("REPORT_INTERVAL")
	if boolReport {
		envTimeReport, err := time.ParseDuration(envStrReport)
		if err != nil {
			log.Println("invalid report interval time in env export")
		} else {
			cfg.PollInterval = envTimeReport
		}
	}

	return cfg
}

func CreateServerConfig() ServerConfig {
	var envcfg Config
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}

	var cfg ServerConfig

	flag.StringVar(&cfg.Address, "a", envcfg.Address, "a string")
	flag.StringVar(&cfg.StoreFile, "f", envcfg.StoreFile, "file store a string")
	flag.StringVar(&cfg.Key, "k", envcfg.Key, "a string")
	flag.StringVar(&cfg.DatabaseDSN, "d", envcfg.DatabaseDSN, "a string")
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

	envStrFile, boolFile := os.LookupEnv("STORE_FILE")
	if boolFile {
		cfg.StoreFile = envStrFile
	}

	envStrDSN, boolDSN := os.LookupEnv("DATABASE_DSN")
	if boolDSN {
		cfg.DatabaseDSN = envStrDSN
	}

	envStrRestore, boolRestore := os.LookupEnv("RESTORE")
	if boolRestore {
		envBoolRestore, err := strconv.ParseBool(envStrRestore)
		if err != nil {
			log.Println("invalid restore bool in env export")
		} else {
			cfg.Restore = envBoolRestore
		}
	}

	envStrStore, boolStore := os.LookupEnv("STORE_INTERVAL")
	if boolStore {
		envTimeStore, err := time.ParseDuration(envStrStore)
		if err != nil {
			log.Println("invalid store interval time in env export")
		} else {
			cfg.StoreInterval = envTimeStore
		}
	}

	return cfg
}
