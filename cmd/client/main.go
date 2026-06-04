package main

import (
	"fmt"
	"time"

	"github.com/dylanoazen/goCache/internal/execution"
)

func main() {
	exec := &execution.Execution{}
	exec.Connect()
	start := time.Now()
	exec.SendMessage("PING")
	procTime, err := exec.ReceiveMessage()
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return
	}
	rtt := time.Since(start)
	fmt.Printf("RTT=%s (procTime=%s)\n", rtt, procTime)
	defer exec.Close()
}
