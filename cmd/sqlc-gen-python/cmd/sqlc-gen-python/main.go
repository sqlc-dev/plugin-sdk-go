package main

import (
	"github.com/tabbed/sqlc-go/codegen"

	python "github.com/tabbed/sqlc-go/cmd/sqlc-gen-python"
)

func main() {
	codegen.Run(python.Generate)
}
