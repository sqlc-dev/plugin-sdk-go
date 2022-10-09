package main

import (
	"context"

	"github.com/tabbed/sqlc-go/codegen"
)

func main() {
	codegen.Run(generate)
}

func generate(ctx context.Context, req *codegen.Request) (*codegen.Response, error) {
	resp := &codegen.Response{
		Files: []*codegen.File{
			{
				Name:     "output.txt",
				Contents: []byte("hello"),
			},
		},
	}
	return resp, nil
}
