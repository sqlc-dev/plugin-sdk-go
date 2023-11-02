package codegen

import (
	"context"

	pb "github.com/sqlc-dev/sqlc-go/plugin"
	"github.com/sqlc-dev/sqlc-go/rpc"
)

type Handler func(context.Context, *pb.GenerateRequest) (*pb.GenerateResponse, error)

func Run(h Handler) {
	rpc.Handle(&server{handler: h})
}
