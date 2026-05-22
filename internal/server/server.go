package server

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
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
		s.sendConnectionConfirm(conn)
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

	sizeStr = strings.TrimSuffix(sizeStr, ":")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	delim, err := reader.ReadByte()
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	if delim != ':' {
		fmt.Println("Read error: invalid checksum delimiter")
		return
	}

	checksumStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	checksumStr = strings.TrimSuffix(checksumStr, "\n")
	checksumValue, err := strconv.ParseUint(checksumStr, 10, 32)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	expected := crc32.ChecksumIEEE(buf)
	if uint32(checksumValue) != expected {
		fmt.Println("Read error: checksum mismatch")
		return
	}

	message := string(buf)

	fmt.Println("received:", message)
	s.mu.Lock()
	s.lastActivity = time.Now()
	s.mu.Unlock()
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

func (s *Server) sendConnectionConfirm(conn net.Conn) {
	message := "Connection confirmed"
	checksum := crc32.ChecksumIEEE([]byte(message))
	_, err := fmt.Fprintf(conn, "%d:%s:%d\n", len(message), message, checksum)
	if err != nil {
		fmt.Println("Write error:", err)
	}
}
