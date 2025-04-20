package initialize

import (
	"context"
	"github.com/ctfloyd/hazelmere-api/src/pkg/client"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_client"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_config"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_logger"
	"github.com/ctfloyd/hazelmere-worker/src/internal/osrs"
)

func InitializeHiscoreClient(config *hz_config.Config, logger hz_logger.Logger) *osrs.HiscoreClient {
	return osrs.NewHiscoreClient(hz_client.NewHttpClient(
		hz_client.HttpClientConfig{
			Host:           config.ValueOrPanic("clients.hiscore.host"),
			TimeoutMs:      config.IntValueOrPanic("clients.hiscore.timeout"),
			Retries:        config.IntValueOrPanic("clients.hiscore.retries"),
			RetryWaitMs:    config.IntValueOrPanic("clients.hiscore.retryWaitMs"),
			RetryMaxWaitMs: config.IntValueOrPanic("clients.hiscore.retryMaxWaitMs"),
		},
		func(msg string) { logger.Error(context.Background(), msg) },
	))
}

func InitializeHazelmereClient(config *hz_config.Config, logger hz_logger.Logger) *client.Hazelmere {
	hazelmere, err := client.NewHazelmere(
		hz_client.NewHttpClient(
			hz_client.HttpClientConfig{
				Host:           config.ValueOrPanic("clients.hazelmere.host"),
				TimeoutMs:      config.IntValueOrPanic("clients.hazelmere.timeout"),
				Retries:        config.IntValueOrPanic("clients.hazelmere.retries"),
				RetryWaitMs:    config.IntValueOrPanic("clients.hazelmere.retryWaitMs"),
				RetryMaxWaitMs: config.IntValueOrPanic("clients.hazelmere.retryMaxWaitMs"),
			},
			func(msg string) { logger.Error(context.Background(), msg) },
		),
		client.HazelmereConfig{
			Token:              config.ValueOrPanic("clients.hazelmere.token"),
			CallingApplication: "hazelmere-worker",
		},
	)
	if err != nil {
		panic(err)
	}
	return hazelmere
}
