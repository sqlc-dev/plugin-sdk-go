package codegen

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	pb "buf.build/gen/go/sqlc/sqlc/protocolbuffers/go/protos/plugin"
	"google.golang.org/protobuf/proto"
)

type Request = pb.CodeGenRequest

type Response = pb.CodeGenResponse

type Handler func(context.Context, *Request) (*Response, error)

func Run(h Handler) {
	if err := run(h); err != nil {
		fmt.Fprintf(os.Stderr, "error generating output: %s", err)
		os.Exit(2)
	}
}

func run(h Handler) error {
	var req Request
	reqBlob, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(reqBlob, &req); err != nil {
		return err
	}
	resp, err := h(context.Background(), &req)
	if err != nil {
		return err
	}
	respBlob, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(os.Stdout)
	if _, err := w.Write(respBlob); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}
