package main

import (
	"math/rand"

	pb "github.com/dougfort/ipdnode/protobuf"
)

func randomMove() pb.MoveMessage_Move {
	choices := []pb.MoveMessage_Move{
		pb.MoveMessage_COOPERATE,
		pb.MoveMessage_DEFECT,
	}

	return choices[rand.Int()%len(choices)]
}
