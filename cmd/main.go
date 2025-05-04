// Package main is a one-file package. Here's only main function.
package main

import (
	// "github.com/mkrtychanr/rag_bot/cmd/commands"

	"github.com/mkrtychanr/rag_bot/cmd/commands"
	_ "go.uber.org/automaxprocs"
)

func main() {
	commands.Execute()
}
