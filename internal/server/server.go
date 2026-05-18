package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	mu           sync.Mutex
	lastActivity time.Time
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	fmt.Println("Server running on :8080")

	s.mu.Lock()
	s.lastActivity = time.Now()
	s.mu.Unlock()

	go s.shutDownIfIdle(ln)

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Println("Accept error:", err)
			continue
		}

		s.mu.Lock()
		s.lastActivity = time.Now()
		s.mu.Unlock()
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Println("received:", message)
	s.mu.Lock()
	s.lastActivity = time.Now()
	s.mu.Unlock()
}

func (s *Server) shutDownIfIdle(ln net.Listener) {
	for {
		time.Sleep(1 * time.Second)
		s.mu.Lock()
		idleFor := time.Since(s.lastActivity)
		s.mu.Unlock()
		if idleFor > 5*time.Second {
			fmt.Println("No activity for 5 Seconds, shutting down server.")
			_ = ln.Close()
			return
		}
	}
}
