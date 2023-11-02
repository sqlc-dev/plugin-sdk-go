package codegen

import (
	"context"

	pb "github.com/sqlc-dev/plugin-sdk-go/plugin"
)

type server struct {
	pb.UnimplementedCodegenServiceServer

	handler Handler
}

func (s *server) Generate(ctx context.Context, req *pb.GenerateRequest) (*pb.GenerateResponse, error) {
	return s.handler(ctx, req)
}
