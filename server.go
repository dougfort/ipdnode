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
	const moveCount = 10
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

	log.Printf("Game: %s: GameStream starts game: %d moves", gameID, moveCount)
	for i := 0; i < moveCount; i++ {
		var theirMove *pb.MoveMessage

		ourMove := pb.MoveMessage{GameID: gameID, Move: randomMove()}

		// send our move to the client
		if err = streamServer.Send(&ourMove); err != nil {
			return errors.Wrapf(err, "Game: %s, Iteration: %d: Send", gameID, i)
		}

		if theirMove, err = streamServer.Recv(); err != nil {
			return errors.Wrapf(err, "Game: %s, Iteration: %d: Recv", gameID, i)
		}

		if theirMove.Move != pb.MoveMessage_COOPERATE && theirMove.Move != pb.MoveMessage_DEFECT {
			return errors.Wrapf(err, "Game: %s: Iteration: %d: Unknown Move: %d",
				gameID, i, theirMove.Move)
		}

		log.Printf("Game: %s; I: %d; Our Move: %s, Their Move: %s",
			gameID,
			i,
			pb.MoveMessage_Move_name[int32(ourMove.Move)],
			pb.MoveMessage_Move_name[int32(theirMove.Move)],
		)
	}

	return nil
}
