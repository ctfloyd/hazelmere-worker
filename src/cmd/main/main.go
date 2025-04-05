package main

import (
	"context"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_config"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_logger"
	"github.com/ctfloyd/hazelmere-worker/src/internal"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := hz_config.NewConfigWithAutomaticDetection()
	err := config.Read()
	if err != nil {
		panic(err)
	}

	logger := hz_logger.NewZeroLogAdapater(hz_logger.LogLevelFromString(config.ValueOrPanic("log.level")))

	app := internal.Application{}
	app.Init(config, logger)
	app.Run(ctx, logger)
	app.Cleanup()
}
