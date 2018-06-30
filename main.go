package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	oldcontext "golang.org/x/net/context"
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
	var grpcClientConn *grpc.ClientConn
	var err error

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if cfg, err = loadConfig(); err != nil {
		log.Printf("loadConfig failed: %v", err)
		return 1
	}

	if listener, err = net.Listen("tcp", cfg.ListenAddress); err != nil {
		log.Printf("net.Listen failed: %v", err)
		return 1
	}

	grpcServer = grpc.NewServer()

	pb.RegisterIPDNodeServer(grpcServer, NewServer())

	go func() {
		log.Printf("Server starts: listening on %s", cfg.ListenAddress)
		if err = grpcServer.Serve(listener); err != nil {
			log.Printf("grpcServer.Serve ended with %s", err)
		}
	}()

	if grpcClientConn, err = grpc.Dial(cfg.DialAddresses[0], grpc.WithInsecure()); err != nil {
		log.Printf("grpc.Dial(%s) failed: %v", cfg.DialAddresses[0], err)
		return 1
	}

	client := pb.NewIPDNodeClient(grpcClientConn)

	gsc, err := client.GameStream(oldcontext.Background())
	if err != nil {
		log.Printf("client.GameStream failed: %v", err)
		return 1
	}

	log.Printf("sending START")
	msg := pb.MoveMessage{GameID: "g.i.d.01", Move: pb.MoveMessage_START}
	if err = gsc.Send(&msg); err != nil {
		log.Printf("gsc.Send failed: %v", err)
		return 1
	}

	s := <-sigChan
	log.Printf("signal: %v; shutting down", s.String())

	grpcServer.GracefulStop()

	return 0
}
