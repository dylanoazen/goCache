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
	_, err := fmt.Fprintf(e.conn, message)
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func (e *Execution) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}