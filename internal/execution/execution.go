package execution

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
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
	messageLength := len(message)
	if messageLength == 0 {
		fmt.Println("Cannot send empty message")
		return
	}

	_, err := fmt.Fprintf(e.conn, "%d:%s", messageLength, message)
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
	resp := make([]byte, 9)
	_, err := io.ReadFull(e.conn, resp)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return
	}

	flag := resp[0]
	procNs := int64(binary.BigEndian.Uint64(resp[1:]))
	procTime := time.Duration(procNs)
	fmt.Printf("Received pong: flag=%d, procTime=%s\n", flag, procTime)
}
