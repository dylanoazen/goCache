package main

import (
	"fmt"
	"github.com/dylanoazen/execution"
	"github.com/dylanoazen/tcp/server"
)

func main() { 
	conn, err := server.Start()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}