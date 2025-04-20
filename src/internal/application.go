package internal

import (
	"context"
	"fmt"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_config"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_logger"
	"github.com/ctfloyd/hazelmere-worker/src/internal/common/handler"
	"github.com/ctfloyd/hazelmere-worker/src/internal/initialize"
	"github.com/ctfloyd/hazelmere-worker/src/internal/snapshot"
	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron/v2"
	"net/http"
)

type Application struct {
	Router    *chi.Mux
	Scheduler gocron.Scheduler
}

func (app *Application) Init(config *hz_config.Config, logger hz_logger.Logger) {
	logger.Info(context.TODO(), "Init Hazelmere worker.")

	hiscoreClient := initialize.InitializeHiscoreClient(config, logger)
	hazelmereClient := initialize.InitializeHazelmereClient(config, logger)

	snapshotService := snapshot.NewSnapshotService(logger, hiscoreClient, hazelmereClient)
	snapshotUpdater := snapshot.NewSnapshotUpdaterJob(logger, hazelmereClient, snapshotService)

	scheduler := initialize.InitializeScheduler(snapshotUpdater)
	app.Scheduler = scheduler

	router := initialize.InitRouter(logger)
	app.Router = router

	snapshotHandler := snapshot.NewSnapshotHandler(logger, snapshotService)
	handlers := []handler.WorkerHandler{
		snapshotHandler,
	}
	for _, h := range handlers {
		h.RegisterRoutes(router, handler.ApiVersionV1)
	}

	logger.Info(context.TODO(), "Done init.")
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
