package syslogish

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	bindHost  = "0.0.0.0"
	bindPort  = 1514
	queueSize = 500
)

// Server implements a UDP-based "syslog-like" server.  Like syslog, as described by RFC 3164, it
// expects that each packet contains a single log message and that, conversely, log messages are
// encapsulated in their entirety by a single packet, however, no attempt is made to parse the
// messages received or validate that they conform to the specification.
type Server struct {
	conn         net.PacketConn
	listening    bool
	storageQueue chan string
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

	return &Server{
		conn:         c,
		storageQueue: make(chan string, queueSize),
	}, nil
}

// Listen starts the server's main loop.
func (s *Server) Listen() {
	// Should only ever be called once
	if !s.listening {
		s.listening = true
		go s.receive()
		go s.parse()
		log.Println("syslogish server running")
	}
}

func (s *Server) receive() {
	// Make buffer the same size as the max for a UDP packet
	buf := make([]byte, 65535)
	for {
		n, _, err := s.conn.ReadFrom(buf)
		if err != nil {
			log.Fatal("syslogish server read error", err)
		}
		message := strings.TrimSuffix(string(buf[:n]), "\n")
		select {
		case s.storageQueue <- message:
		default:
		}
	}
}

func (s *Server) parse() {
	for message := range s.storageQueue {
		curlyIndex := strings.Index(message, "{")
		if curlyIndex > -1 {
			message = message[curlyIndex:]
			var messageJSON map[string]interface{}
			err := json.Unmarshal([]byte(message), &messageJSON)
			if err == nil && messageJSON["kubernetes"] != nil {
				log.Printf("NGINX LOG:%+v\n\n", messageJSON)
			}
		}
	}
}
