package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {
	err := app.RunContext(initCtx(), os.Args)
	if err != nil {
		panic(err)
	}
}

func initCtx() context.Context {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sig
		cancel()
		// don't close channel, caus it may cause panic, when app receive
		// multiple signals, and writing them to closed channel
	}()
	return ctx
}

var app = &cli.App{
	Name:        "covercut",
	Usage:       "🔪 Cut your cover profiles Jack the Ripper! 👻",
	Action:      Action,
	Description: "Parses tl schema and generates go code for it",
	Commands:    []*cli.Command{},
	Flags:       []cli.Flag{},
	Copyright:   "Xelaj Software",

	Reader:    os.Stdin,
	Writer:    os.Stdout,
	ErrWriter: os.Stderr,

	HideHelpCommand: true,
}

type Exiter struct {
	Code int
}

func (e Exiter) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("Code %v", e.Code)
	}
	return "Just for test"
}

func (e Exiter) ExitCode() int { return e.Code }
