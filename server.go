package main

import (
	"log"

	"github.com/pkg/errors"

	pb "github.com/dougfort/ipdnode/protobuf"
)

type serverState struct {
}

// NewServer returns an object that implements  pb.IPDNodeServer
func NewServer() pb.IPDNodeServer {
	return &serverState{}
}

// GameStream is the method of pb.IPDNodeServer
func (s *serverState) GameStream(streamServer pb.IPDNode_GameStreamServer) error {
	var message *pb.MoveMessage
	var gameID string
	var err error

	// we expect to get a start to start with
	if message, err = streamServer.Recv(); err != nil {
		return errors.Wrap(err, "init Recv")
	}

	if message.GetMove() != pb.MoveMessage_START {
		return errors.Errorf("%s: unexpected move, expecting start of game", message.GetGameID())
	}

	if gameID = message.GetGameID(); gameID == "" {
		return errors.Errorf("invalid empty gameID")
	}

	log.Printf("%s: GameStream starts", gameID)
	return nil

}
