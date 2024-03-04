package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	runAddr    string
	dbURI      string
	accSysAddr string
}

func (sa *ServerConfig) GetRunAddress() string {
	return sa.runAddr
}

func (sa *ServerConfig) GetDbURI() string {
	return sa.dbURI
}

func (sa *ServerConfig) GetAccSysAddr() string {
	return sa.accSysAddr
}

func (sa *ServerConfig) SetRunAddress(value string) {
	sa.runAddr = value
}

func (sa *ServerConfig) SetDbURI(value string) {
	sa.dbURI = value
}

func (sa *ServerConfig) SetAccSysAddr(value string) {
	sa.accSysAddr = value
}

func (sa *ServerConfig) ParseFlags() {
	flag.StringVar(&sa.runAddr, "a", "localhost:8000", "address and port to run shortener")
	flag.StringVar(&sa.dbURI, "d", "", "db address")
	flag.StringVar(&sa.accSysAddr, "r", "http://localhost:8080", "accrual system address")

	flag.Parse()

	if envRunAddr, in := os.LookupEnv("RUN_ADDRESS"); in {
		sa.SetRunAddress(envRunAddr)
	}

	if envDbURI, in := os.LookupEnv("DATABASE_URI"); in {
		sa.SetDbURI(envDbURI)
	}

	if envAccSysAddr, in := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); in {
		sa.SetAccSysAddr(envAccSysAddr)
	}

}
