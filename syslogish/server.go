package syslogish

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/deis/stdout-metrics/influx"
	"github.com/deis/stdout-metrics/util"
	"github.com/influxdata/influxdb/client/v2"
)

const (
	bindHost  = "0.0.0.0"
	bindPort  = 1514
	queueSize = 50000
)

// Server implements a UDP-based "syslog-like" server.  Like syslog, as described by RFC 3164, it
// expects that each packet contains a single log message and that, conversely, log messages are
// encapsulated in their entirety by a single packet, however, no attempt is made to parse the
// messages received or validate that they conform to the specification.
type Server struct {
	Conn         net.PacketConn
	Listening    bool
	StorageQueue chan string
	InfluxClient client.Client
}

// NewServer returns a pointer to a new Server instance.
func NewServer() (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", bindHost, bindPort))
	if err != nil {
		return nil, err
	}
	c, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	influxClient, err := influx.Connect()
	if err != nil {
		return nil, err
	}

	return &Server{
		Conn:         c,
		StorageQueue: make(chan string, queueSize),
		InfluxClient: influxClient,
	}, nil
}

// Listen starts the server's main loop.
func (s *Server) Listen() {
	// Should only ever be called once
	if !s.Listening {
		defer s.InfluxClient.Close()
		s.Listening = true
		go s.receive()
		go s.parse()
		log.Println("syslogish server running")
	}
}

func (s *Server) receive() {
	// Make buffer the same size as the max for a UDP packet
	buf := make([]byte, 65535)
	for {
		n, _, err := s.Conn.ReadFrom(buf)
		if err != nil {
			log.Fatal("syslogish server read error", err)
		}
		message := strings.TrimSuffix(string(buf[:n]), "\n")
		select {
		case s.StorageQueue <- message:
		default:
		}
	}
}

func (s *Server) parse() {
	for message := range s.StorageQueue {
		messageJSON, err := util.ParseMessage(message)
		if err == nil && util.FromContainer(messageJSON, "deis-router") {
			err = influx.WriteRouterMetrics(s.InfluxClient, messageJSON)
			if err != nil {
				fmt.Printf("Error@Metrics:%v\n", err)
			}
		}
	}
}
