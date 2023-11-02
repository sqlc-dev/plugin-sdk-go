package codegen

import (
	"context"

	"github.com/sqlc-dev/sqlc-go/internal/rpc"
	pb "github.com/sqlc-dev/sqlc-go/plugin"
)

type Handler func(context.Context, *pb.GenerateRequest) (*pb.GenerateResponse, error)

func Run(h Handler) {
	rpc.Handle(&server{handler: h})
}
