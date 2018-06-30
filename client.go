package main

import (
	"log"

	"github.com/pkg/errors"

	oldcontext "golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/dougfort/ipdnode/protobuf"
)

func runClient(clientNum int, grpcClientConn *grpc.ClientConn) error {
	var err error

	client := pb.NewIPDNodeClient(grpcClientConn)

	gsc, err := client.GameStream(oldcontext.Background())
	if err != nil {
		return errors.Wrap(err, "client.GameStream failed")
	}

	log.Printf("client #%d sending START", clientNum)
	msg := pb.MoveMessage{GameID: "g.i.d.01", Move: pb.MoveMessage_START}
	if err = gsc.Send(&msg); err != nil {
		return errors.Wrap(err, "gsc.Send failed")
	}

	return nil
}
