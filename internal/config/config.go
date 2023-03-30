package config

import (
	"flag"
	"github.com/caarlos0/env/v7"
	"log"
	"os"
	"strconv"
	"time"
)

type EnvConfig struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreFile      string        `env:"STORE_FILE" envDefault:"tmp/devops-metrics-db.json"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	Restore        bool          `env:"RESTORE" envDefault:"true"`
}

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func CreateAgentConfig() AgentConfig {
	var envcfg EnvConfig
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}

	var conf AgentConfig

	flagStrAddr := flag.String("a", "127.0.0.1:8080", "a string")
	flagStrPollInterval := flag.Duration("p", 2*time.Second, "a duration")
	flagStrReportInterval := flag.Duration("r", 10*time.Second, "a duration")
	flag.Parse()

	envStrAddres, bool := os.LookupEnv("ADDRESS")
	if bool {
		conf.Address = envStrAddres
	} else {
		conf.Address = *flagStrAddr
	}

	envStrPollInterval, bool := os.LookupEnv("POLL_INTERVAL")
	if bool {
		envTimePollInterval, err := time.ParseDuration(envStrPollInterval)
		if err != nil {
			conf.PollInterval = envcfg.PollInterval
			log.Println("incorrect time value restore in env export")
		} else {
			conf.PollInterval = envTimePollInterval
		}
	} else {
		conf.PollInterval = *flagStrPollInterval
	}

	envStrReportInterval, bool := os.LookupEnv("REPORT_INTERVAL")
	if bool {
		envTimeReportInterval, err := time.ParseDuration(envStrReportInterval)
		if err != nil {
			conf.ReportInterval = envcfg.ReportInterval
			log.Println("incorrect time value restore in env export")
		} else {
			conf.PollInterval = envTimeReportInterval
		}
	} else {
		conf.ReportInterval = *flagStrReportInterval
	}

	return conf
}

type ServerConfig struct {
	Address       string
	StoreFile     string
	Restore       bool
	StoreInterval time.Duration
}

func CreateServerConfig() ServerConfig {
	var envcfg EnvConfig
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}

	var conf ServerConfig

	flagStrAddr := flag.String("a", "127.0.0.1:8080", "a string")
	flagStrFile := flag.String("f", "tmp/devops-metrics-db.json", "a string")
	flagStrRestore := flag.Bool("r", true, "a bool")
	flagStrStoreInterval := flag.Duration("i", 300*time.Second, "a duration")
	flag.Parse()

	envStrAddres, bool := os.LookupEnv("ADDRESS")
	if bool {
		conf.Address = envStrAddres
	} else {
		conf.Address = *flagStrAddr
	}

	evnStrFile, bool := os.LookupEnv("STORE_FILE")
	if bool {
		conf.StoreFile = evnStrFile
	} else {
		conf.StoreFile = *flagStrFile
	}

	envStrRestore, bool := os.LookupEnv("RESTORE")
	if bool {
		envBoolRestore, err := strconv.ParseBool(envStrRestore)
		if err != nil {
			conf.Restore = envcfg.Restore
			log.Println("incorrect bool value restore in env export")
		} else {
			conf.Restore = envBoolRestore
		}
	} else {
		conf.Restore = *flagStrRestore
	}

	envStrStoreInterval, bool := os.LookupEnv("STORE_INTERVAL")
	if bool {
		envTimeStoreInterval, err := time.ParseDuration(envStrStoreInterval)
		if err != nil {
			conf.StoreInterval = envcfg.StoreInterval
			log.Println("incorrect time value restore in env export")
		} else {
			conf.StoreInterval = envTimeStoreInterval
		}
	} else {
		conf.StoreInterval = *flagStrStoreInterval
	}

	return conf
}
