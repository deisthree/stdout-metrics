package main

import (
	"github.com/deis/nginxpusher/syslogish"
	"log"
	"os"
)

func main() {
	syslogishServer, err := syslogish.NewServer()
	if err != nil {
		log.Fatal("Error creating syslogish server", err)
	}
	syslogishServer.Listen()

	log.Println("Nginx Pusher running")
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}
