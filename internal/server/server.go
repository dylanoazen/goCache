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

	"github.com/dylanoazen/goCache/internal/cache"
)

type Server struct {
	mu           sync.Mutex
	lastActivity time.Time
	cacheMu      sync.RWMutex
	cache        *cache.Cache
}

func (s *Server) Start() {
	if s.cache == nil {
		s.cache = cache.NewCache()
	}
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	fmt.Println("Server running on :8080")

	s.updateLastActivity()

	go s.shutDownIfIdle(ln)
	s.serve(ln)
}

func (s *Server) serve(ln net.Listener) {

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Println("Accept error:", err)
			continue
		}
		s.updateLastActivity()
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	message, recvTime, err := s.readMessage(conn)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Println("received:", message)
	s.updateLastActivity()
	s.handleMessage(conn, message, recvTime)
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

func (s *Server) readMessage(conn net.Conn) (string, time.Time, error) {
	reader := bufio.NewReader(conn)
	sizeStr, err := reader.ReadString(':')
	if err != nil {
		return "", time.Time{}, err
	}
	recvTime := time.Now()
	sizeStr = strings.TrimSuffix(sizeStr, ":")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return "", time.Time{}, err
	}
	buf := make([]byte, size)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return "", time.Time{}, err
	}
	return string(buf), recvTime, nil
}

func (s *Server) handleMessage(conn net.Conn, message string, recvTime time.Time) {
	parts := strings.Fields(message)
	if len(parts) == 0 {
		return
	}

	switch strings.ToUpper(parts[0]) {
	case "PING":
		elapsed := time.Since(recvTime)
		s.sendPong(conn, elapsed)
	case "SET":
		if len(parts) < 3 {
			return
		}
		key := parts[1]
		value := strings.Join(parts[2:], " ")
		s.cacheMu.Lock()
		s.cache.Set(key, value)
		s.cacheMu.Unlock()
	case "GET":
		if len(parts) < 2 {
			return
		}
		key := parts[1]
		s.cacheMu.RLock()
		value, ok := s.cache.Get(key)
		s.cacheMu.RUnlock()
		if !ok {
			s.sendTextResponse(conn, "(nil)")
			return
		}
		s.sendTextResponse(conn, value)
	}
}

func (s *Server) updateLastActivity() {
	s.mu.Lock()
	s.lastActivity = time.Now()
	s.mu.Unlock()
}

func (s *Server) sendTextResponse(conn net.Conn, message string) {
	messageLength := len(message)
	if messageLength == 0 {
		message = ""
		messageLength = 0
	}
	_, err := fmt.Fprintf(conn, "%d:%s", messageLength, message)
	if err != nil {
		fmt.Println("Write error:", err)
	}
}
