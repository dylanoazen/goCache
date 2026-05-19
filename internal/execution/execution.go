package execution

import (
	"fmt"
	"net"
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
	buf := make([]byte, 1024)
	n, err := e.conn.Read(buf)
	if err != nil {
		fmt.Println("Error receiving message:", err)
		return
	}
	fmt.Println("Received from server:", string(buf[:n]))
}