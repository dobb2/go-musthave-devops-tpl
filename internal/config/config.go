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
	valAddr, boolAdrr := os.LookupEnv("ADDRESS")
	if boolAdrr {
		flag.StringVar(&conf.Address, "a", valAddr, "a string")
	} else {
		flag.StringVar(&conf.Address, "a", envcfg.Address, "a string")
	}

	strPol, boolPol := os.LookupEnv("POLL_INTERVAL")
	if boolPol {
		valPol, err := time.ParseDuration(strPol)
		if err == nil {
			flag.DurationVar(&conf.PollInterval, "p", valPol, "a durations")
		} else {
			log.Println(err)
			flag.DurationVar(&conf.PollInterval, "p", envcfg.PollInterval, "a duration")
		}
	} else {
		flag.DurationVar(&conf.PollInterval, "p", envcfg.PollInterval, "a duration")
	}

	strRep, boolRep := os.LookupEnv("REPORT_INTERVAL")
	if boolRep {
		valRep, err := time.ParseDuration(strRep)
		if err == nil {
			flag.DurationVar(&conf.ReportInterval, "r", valRep, "a durations")
		} else {
			log.Println(err)
			log.Println("this problem 11")
			flag.DurationVar(&conf.ReportInterval, "r", envcfg.ReportInterval, "a duration")
		}
	} else {
		flag.DurationVar(&conf.ReportInterval, "r", envcfg.ReportInterval, "a duration")
	}

	flag.Parse()
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

	str, bool := os.LookupEnv("ADDRESS")
	if bool {
		flag.StringVar(&conf.Address, "a", str, "a string")
	} else {
		flag.StringVar(&conf.Address, "a", envcfg.Address, "a string")
	}

	str, bool = os.LookupEnv("STORE_INTERVAL")
	if bool {
		val, err := time.ParseDuration(str)
		if err == nil {
			flag.DurationVar(&conf.StoreInterval, "i", val, "a durations")
		} else {
			log.Println(err)
			flag.DurationVar(&conf.StoreInterval, "i", envcfg.StoreInterval, "a duration")
		}
	} else {
		flag.DurationVar(&conf.StoreInterval, "i", envcfg.StoreInterval, "a duration")
	}

	str, bool = os.LookupEnv("STORE_FILE")
	if bool {
		flag.StringVar(&conf.StoreFile, "f", str, "a string")
	} else {
		flag.StringVar(&conf.StoreFile, "f", envcfg.StoreFile, "a string")
	}

	str, bool = os.LookupEnv("RESTORE")
	if bool {
		valBool, err := strconv.ParseBool(str)
		if err == nil {
			flag.BoolVar(&conf.Restore, "r", valBool, "a bool")
		} else {
			log.Println(err)
			flag.BoolVar(&conf.Restore, "r", envcfg.Restore, "a bool")
		}
	} else {
		flag.BoolVar(&conf.Restore, "r", envcfg.Restore, "a bool")
	}
	flag.Parse()
	return conf
}
