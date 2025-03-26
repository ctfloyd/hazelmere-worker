package job

import (
	"api/src/internal/osrs"
	"context"
	"github.com/ctfloyd/hazelmere-api/src/pkg/api"
	"github.com/ctfloyd/hazelmere-api/src/pkg/client"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_logger"
	"regexp"
	"slices"
	"strings"
	"time"
)

type SnapshotUpdaterJob struct {
	hiscoreClient   *osrs.HiscoreClient
	hazelmereClient *client.Hazelmere
	logger          hz_logger.Logger
}

func NewSnapshotUpdaterJob(hiscoreClient *osrs.HiscoreClient, hazelmereClient *client.Hazelmere, logger hz_logger.Logger) *SnapshotUpdaterJob {
	return &SnapshotUpdaterJob{
		hiscoreClient:   hiscoreClient,
		hazelmereClient: hazelmereClient,
		logger:          logger,
	}
}

func (job SnapshotUpdaterJob) Run() {
	ctx := context.Background()

	job.logger.Info(ctx, "Running the snapshot updater job.")

	usersResponse, err := job.hazelmereClient.User.GetAllUsers()
	if err != nil {
		job.logger.Error(ctx, "An error occurred while getting all users: "+err.Error())
		return
	}

	for _, user := range usersResponse.Users {
		if user.TrackingStatus == api.TrackingStatusEnabled {
			job.logger.InfoArgs(ctx, "Generating snapshot for user: "+user.RunescapeName)

			hiscore, err := job.hiscoreClient.GetHiscore(user.RunescapeName)
			if err != nil {
				job.logger.Error(ctx, "An error occurred while getting hiscore: "+err.Error())
				continue
			}

			snapshot := makeSnapshot(user, hiscore)

			snapshotRequest := api.CreateSnapshotRequest{
				Snapshot: snapshot,
			}
			_, err = job.hazelmereClient.Snapshot.CreateSnapshot(snapshotRequest)
			if err != nil {
				job.logger.Error(ctx, "An error occurred while creating snapshot: "+err.Error())
				continue
			}

			job.logger.Info(ctx, "Snapshot created successfully!")
		} else {
			job.logger.InfoArgs(ctx, "Skipping generating snapshot for user: "+user.RunescapeName+" because tracking is disabled.")
		}
	}
}

func makeSnapshot(user api.User, hiscore osrs.Hiscore) api.HiscoreSnapshot {
	snapshot := api.HiscoreSnapshot{
		UserId:     user.Id,
		Timestamp:  time.Now(),
		Skills:     make([]api.SkillSnapshot, 0),
		Bosses:     make([]api.BossSnapshot, 0),
		Activities: make([]api.ActivitySnapshot, 0),
	}

	for _, skill := range hiscore.Skills {
		snapshot.Skills = append(snapshot.Skills, api.SkillSnapshot{
			ActivityType: api.ActivityTypeFromValue(convertToActivityTypeEnumValue(skill.Name)),
			Name:         skill.Name,
			Level:        skill.Level,
			Experience:   skill.Xp,
			Rank:         skill.Rank,
		})
	}

	for _, activity := range hiscore.Activities {
		at := api.ActivityTypeFromValue(convertToActivityTypeEnumValue(activity.Name))

		if slices.Contains(api.AllBossActivityTypes, at) {
			snapshot.Bosses = append(snapshot.Bosses, api.BossSnapshot{
				ActivityType: at,
				Name:         activity.Name,
				KillCount:    activity.Score,
				Rank:         activity.Rank,
			})
		} else {
			snapshot.Activities = append(snapshot.Activities, api.ActivitySnapshot{
				ActivityType: at,
				Name:         activity.Name,
				Score:        activity.Score,
				Rank:         activity.Rank,
			})
		}
	}

	return snapshot
}

func removeIllegalEnumCharacters(str string) string {
	remove := regexp.MustCompile("['\\-_:()]")
	return remove.ReplaceAllString(str, "")
}

func convertToActivityTypeEnumValue(name string) string {
	name = removeIllegalEnumCharacters(name)
	return strings.ToUpper(strings.ReplaceAll(name, " ", "_"))
}
