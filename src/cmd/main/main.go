package main

import (
	"api/src/internal"
	"api/src/internal/common/logger"
	"context"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	l := logger.NewZeroLogAdapater(logger.LogLevelDebug)

	app := internal.Application{}
	app.Init(l)
	app.Run(ctx, l)
	app.Cleanup()
}
