package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/deis/stdout-metrics/syslogish"
)

func main() {

	syslogishServer, err := syslogish.NewServer()
	if err != nil {
		log.Fatal("Error creating syslogish server", err)
	}
	syslogishServer.Listen()
	log.Println("stdout metrics is running!")

	reopen := make(chan os.Signal, 1)
	signal.Notify(reopen, syscall.SIGTERM)

	for {
		<-reopen
		log.Println("Clean up all the stuffs")
	}
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}
