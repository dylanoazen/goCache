package execution

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"strconv"
	"strings"
)

type Execution struct {
	conn net.Conn
}

func (e *Execution) Connect() {         
	coon, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	e.conn = coon
}

func (e *Execution) SendMessage(message string) {
	if e.conn == nil {
		fmt.Println("No connection available")
		return
	}

	messageLength := len(message)
	if messageLength == 0 {
		fmt.Println("Cannot send empty message")
		return
	}

	checksum := crc32.ChecksumIEEE([]byte(message))
	_, err := fmt.Fprintf(e.conn, "%d:%s:%d\n", messageLength, message, checksum)
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func (e *Execution) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

func (e *Execution) ReceiveMessage() {
	if e.conn == nil {
		fmt.Println("No connection available")
		return
	}

	reader := bufio.NewReader(e.conn)
	sizeStr, err := reader.ReadString(':')
	if err != nil {
		fmt.Println("Error receiving message:", err)
		return
	}

	sizeStr = strings.TrimSuffix(sizeStr, ":")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		fmt.Println("Error receiving message:", err)
		return
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		fmt.Println("Error receiving message:", err)
		return
	}

	delim, err := reader.ReadByte()
	if err != nil {
		fmt.Println("Error receiving message:", err)
		return
	}
	if delim != ':' {
		fmt.Println("Error receiving message: invalid checksum delimiter")
		return
	}

	checksumStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error receiving message:", err)
		return
	}

	checksumStr = strings.TrimSuffix(checksumStr, "\n")
	checksumValue, err := strconv.ParseUint(checksumStr, 10, 32)
	if err != nil {
		fmt.Println("Error receiving message:", err)
		return
	}

	expected := crc32.ChecksumIEEE(buf)
	if uint32(checksumValue) != expected {
		fmt.Println("Error receiving message: checksum mismatch")
		return
	}

	fmt.Println("Received from server:", string(buf))
}