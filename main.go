package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	os.Exit(run())
}

func run() int {
	const serverAddress = "localhost:1111"

	var listener net.Listener
	var grpcServer *grpc.Server
	var err error

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if listener, err = net.Listen("tcp", serverAddress); err != nil {
		log.Printf("net.Listen failled: %v", err)
		return 1
	}

	grpcServer = grpc.NewServer()

	go func() {
		log.Printf("Server starts: listening on %s", serverAddress)
		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("grpcServer.Serve ended with %s", err)
		}
	}()

	s := <-sigChan
	log.Printf("signal: %v; shutting down", s.String())

	grpcServer.GracefulStop()

	return 0
}
