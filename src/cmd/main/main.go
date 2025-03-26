package main

import (
	"api/src/internal"
	"context"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_logger"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	l := hz_logger.NewZeroLogAdapater(hz_logger.LogLevelDebug)

	app := internal.Application{}
	app.Init(l)
	app.Run(ctx, l)
	app.Cleanup()
}
