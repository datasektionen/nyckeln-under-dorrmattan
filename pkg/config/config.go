package config

import (
	"flag"
	"os"
)

type Config struct {
	PlsPort    string
	LoginPort  string
	SsoPort    string
	HivePort   string
	HodisURL   string
	InitKTHID  string
	ConfigFile string
}

var config Config

func init() {
	flag.StringVar(&config.PlsPort, "pls-port", "7001", "Port for the pls service")
	flag.StringVar(&config.LoginPort, "login-port", "7002", "Port for the login service")
	flag.StringVar(&config.SsoPort, "sso-port", "7003", "Port for the sso service")
	flag.StringVar(&config.HivePort, "hive-port", "7004", "Port for the hive service")
	flag.StringVar(&config.HodisURL, "hodis-url", "https://hodis.datasektionen.se", "URL to the hodis instance")
	flag.StringVar(&config.InitKTHID, "kth-id", os.Getenv("KTH_ID"), "Username to use for login")
	flag.StringVar(&config.ConfigFile, "config-file", "config.yaml", "Path to a yaml config file")
}

func GetConfig() *Config {
	flag.Parse()

	return &config
}
