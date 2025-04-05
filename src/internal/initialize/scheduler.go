package initialize

import (
	"api/src/internal/snapshot"
	"github.com/go-co-op/gocron/v2"
	"time"
)

func InitializeScheduler(snapshotUpdater *snapshot.SnapshotUpdaterJob) gocron.Scheduler {
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

	return scheduler
}
