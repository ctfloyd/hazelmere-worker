package snapshot

import (
	"context"
	"github.com/ctfloyd/hazelmere-api/src/pkg/api"
	"github.com/ctfloyd/hazelmere-api/src/pkg/client"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_logger"
)

type SnapshotUpdaterJob struct {
	logger          hz_logger.Logger
	hazelmere       *client.Hazelmere
	snapshotService SnapshotService
}

func NewSnapshotUpdaterJob(logger hz_logger.Logger, hazelmere *client.Hazelmere, snapshotService SnapshotService) *SnapshotUpdaterJob {
	return &SnapshotUpdaterJob{
		logger:          logger,
		hazelmere:       hazelmere,
		snapshotService: snapshotService,
	}
}

func (job SnapshotUpdaterJob) Run() {
	ctx := context.Background()

	job.logger.Info(ctx, "Running the snapshot updater job.")

	usersResponse, err := job.hazelmere.User.GetAllUsers()
	if err != nil {
		job.logger.ErrorArgs(ctx, "An error occurred while getting all users: %v", err)
		return
	}

	for _, user := range usersResponse.Users {
		if user.TrackingStatus == api.TrackingStatusEnabled {
			job.logger.InfoArgs(ctx, "Generating snapshot for user: %s", user.RunescapeName)

			_, err := job.snapshotService.MakeSnapshot(user)
			if err != nil {
				job.logger.ErrorArgs(ctx, "Failed to create snapshot for user %v: %v.", user, err)
			}

			job.logger.Info(ctx, "Snapshot created successfully!")
		} else {
			job.logger.InfoArgs(ctx, "Skipping generating snapshot for user: "+user.RunescapeName+" because tracking is disabled.")
		}
	}
}
