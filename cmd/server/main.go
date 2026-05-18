package main

import "github.com/dylanoazen/goCache/internal/server"

func main() {
	srv := &server.Server{}
	srv.Start()
}
