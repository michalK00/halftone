package main

import (
	"context"
	_ "github.com/michalK00/sg-qr/docs"
	"github.com/michalK00/sg-qr/internal/cmd"
	"os"
	"os/signal"
)

// @title Halftone
// @version 0.1
// @contact.name Michał Klemens
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	ret := cmd.Execute(ctx)
	os.Exit(ret)
}
