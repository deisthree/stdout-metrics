package syslogish

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
)

const (
	bindHost  = "0.0.0.0"
	bindPort  = 1514
	queueSize = 500
)

var appRegex *regexp.Regexp

func init() {
	appRegex = regexp.MustCompile(`^.* ([-_a-z0-9]+)\[[a-z0-9-_\.]+\].*`)
}

// Server implements a UDP-based "syslog-like" server.  Like syslog, as described by RFC 3164, it
// expects that each packet contains a single log message and that, conversely, log messages are
// encapsulated in their entirety by a single packet, however, no attempt is made to parse the
// messages received or validate that they conform to the specification.
type Server struct {
	conn         net.PacketConn
	listening    bool
	storageQueue chan string
	adapterMutex sync.RWMutex
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
		// Strip off all the leading syslog junk and just take the JSON.
		// Drop messages that clearly do not contain any JSON, although an open curly brace is only
		// a soft indicator of JSON.  If the message does not contain JSON or is otherwise malformed,
		// it may still be dropped when parsing is attempted.
		curlyIndex := strings.Index(message, "{")
		if curlyIndex > -1 {
			message = message[curlyIndex:]
			// Parse the message into json
			var messageJSON map[string]interface{}
			_ = json.Unmarshal([]byte(message), &messageJSON)
			log.Printf("NGINX LOG:%+v", messageJSON)
			s.adapterMutex.RUnlock()
		}
	}
}
