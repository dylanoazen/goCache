package main

import (
	"github.com/dylanoazen/goCache/internal/execution"
)

func main() {
	exec := &execution.Execution{}
	exec.Connect()
	exec.SendMessage("PING")
	exec.ReceiveMessage()
	defer exec.Close()
}
