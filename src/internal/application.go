package internal

import (
	"api/src/internal/initialize"
	"api/src/internal/job"
	"api/src/internal/osrs"
	"context"
	"fmt"
	"github.com/ctfloyd/hazelmere-api/src/pkg/client"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_client"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron/v2"
	"net/http"
	"time"
)

type Application struct {
	Router    *chi.Mux
	Scheduler gocron.Scheduler
}

func (app *Application) Init(l hz_logger.Logger) {
	l.Info(context.TODO(), "Init Hazelmere worker.")

	router := initialize.InitRouter(l)
	app.Router = router

	hsClient := osrs.NewHiscoreClient(hz_client.NewHttpClient(
		hz_client.HttpClientConfig{
			Host:           "https://secure.runescape.com/m=hiscore_oldschool/index_lite.json",
			TimeoutMs:      20000,
			Retries:        2,
			RetryWaitMs:    50,
			RetryMaxWaitMs: 500,
		},
		func(msg string) { l.Error(context.Background(), msg) },
	))

	hazelmereClient := client.NewHazelmere(
		hz_client.NewHttpClient(
			hz_client.HttpClientConfig{
				Host:           "https://api.hazelmere.xyz",
				TimeoutMs:      1000,
				Retries:        1,
				RetryWaitMs:    50,
				RetryMaxWaitMs: 500,
			},
			func(msg string) { l.Error(context.Background(), msg) },
		),
	)

	snapshotUpdater := job.NewSnapshotUpdaterJob(hsClient, hazelmereClient, l)

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	topOfNextHour := time.Now().Truncate(time.Hour).Add(1 * time.Hour)
	_, err = scheduler.NewJob(
		gocron.DurationJob(60*time.Minute),
		gocron.NewTask(snapshotUpdater.Run),
		gocron.WithStartAt(gocron.WithStartDateTime(topOfNextHour)),
	)
	if err != nil {
		panic(err)
	}
	app.Scheduler = scheduler

	l.Info(context.TODO(), "Done init.")
}

func (app *Application) Run(ctx context.Context, l hz_logger.Logger) {
	app.Scheduler.Start()

	l.Info(ctx, "Trying listen and serve 8080.")
	err := http.ListenAndServe(":8080", app.Router)
	if err != nil {
		panic(err)
	}
}

func (app *Application) Cleanup() {
	fmt.Println("Cleaning up Hazelmere worker.")
	_ = app.Scheduler.Shutdown()
}
