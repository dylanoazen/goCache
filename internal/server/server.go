package server

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
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

	reader := bufio.NewReader(conn)
	sizeStr, err := reader.ReadString(':')
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	recvTime := time.Now()
	sizeStr = strings.TrimSuffix(sizeStr, ":")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		fmt.Println("Parse size error:", err)
		return
	}
	buf := make([]byte, size)
	_, err = io.ReadFull(reader, buf)
	message := string(buf)

	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Println("received:", message)
	s.mu.Lock()
	s.lastActivity = time.Now()
	s.mu.Unlock()

	if strings.EqualFold(message, "PING") {
		elapsed := time.Since(recvTime)
		s.sendPong(conn, elapsed)
	}
}

func (s *Server) shutDownIfIdle(ln net.Listener) {
	for {
		time.Sleep(1 * time.Minute)
		s.mu.Lock()
		idleFor := time.Since(s.lastActivity)
		s.mu.Unlock()
		if idleFor > 10*time.Minute {
			fmt.Println("No activity for 10 Minutes, shutting down server.")
			_ = ln.Close()
			return
		}
	}
}

func (s *Server) sendPong(conn net.Conn, elapsed time.Duration) {
	resp := make([]byte, 9)
	resp[0] = 1
	binary.BigEndian.PutUint64(resp[1:], uint64(elapsed.Nanoseconds()))
	_, err := conn.Write(resp)
	if err != nil {
		fmt.Println("Write error:", err)
	}
}
