package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	pb "github.com/dougfort/ipdnode/protobuf"
)

func main() {
	os.Exit(run())
}

func run() int {
	var cfg ConfigType
	var listener net.Listener
	var grpcServer *grpc.Server
	var err error

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if cfg, err = loadConfig(); err != nil {
		log.Printf("loadConfig failed: %v", err)
		return 1
	}

	if listener, err = net.Listen("tcp", cfg.ListenAddress); err != nil {
		log.Printf("net.Listen failled: %v", err)
		return 1
	}

	grpcServer = grpc.NewServer()

	pb.RegisterIPDNodeServer(grpcServer, NewServer())

	go func() {
		log.Printf("Server starts: listening on %s", cfg.ListenAddress)
		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("grpcServer.Serve ended with %s", err)
		}
	}()

	s := <-sigChan
	log.Printf("signal: %v; shutting down", s.String())

	grpcServer.GracefulStop()

	return 0
}
