package main

import (
	"log"

	"github.com/pkg/errors"

	oldcontext "golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/dougfort/ipdnode/protobuf"
)

func runClient(clientNum int, grpcClientConn *grpc.ClientConn) error {
	const moveCount = 10
	var err error

	client := pb.NewIPDNodeClient(grpcClientConn)

	gsc, err := client.GameStream(oldcontext.Background())
	if err != nil {
		return errors.Wrap(err, "client.GameStream failed")
	}

	gameID := "g.i.d.01"

	log.Printf("client #%d sending START", clientNum)
	message := pb.MoveMessage{GameID: gameID, Move: pb.MoveMessage_START}
	if err = gsc.Send(&message); err != nil {
		return errors.Wrap(err, "gsc.Send failed")
	}

	log.Printf("Game: %s: GameStream starts game: %d moves", gameID, moveCount)
	for i := 0; i < moveCount; i++ {
		var theirMove *pb.MoveMessage

		if theirMove, err = gsc.Recv(); err != nil {
			return errors.Wrapf(err, "Game: %s, Iteration: %d: Recv", gameID, i)
		}

		if theirMove.Move != pb.MoveMessage_COOPERATE && theirMove.Move != pb.MoveMessage_DEFECT {
			return errors.Wrapf(err, "Game: %s: Iteration: %d: Unknown Move: %d",
				gameID, i, theirMove.Move)
		}

		ourMove := pb.MoveMessage{GameID: gameID, Move: randomMove()}

		// send our move to the client
		if err = gsc.Send(&ourMove); err != nil {
			return errors.Wrapf(err, "Game: %s, Iteration: %d: Send", gameID, i)
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
