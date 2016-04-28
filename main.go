package main

import (
	"log"
	"os"

	"github.com/deis/stdout-metrics/syslogish"
)

func main() {
	syslogishServer, err := syslogish.NewServer()
	if err != nil {
		log.Fatal("Error creating syslogish server", err)
	}
	syslogishServer.Listen()
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}
