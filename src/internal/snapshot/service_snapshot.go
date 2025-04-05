package snapshot

import (
	"errors"
	"github.com/ctfloyd/hazelmere-api/src/pkg/api"
	"github.com/ctfloyd/hazelmere-api/src/pkg/client"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_logger"
	"github.com/ctfloyd/hazelmere-worker/src/internal/osrs"
	"regexp"
	"slices"
	"strings"
	"time"
)

var ErrSnapshotService = errors.New("snapshot service error")
var ErrUserNotFound = errors.New("user not found")

type SnapshotService interface {
	MakeSnapshotForUser(userId string) (api.HiscoreSnapshot, error)
	MakeSnapshot(user api.User) (api.HiscoreSnapshot, error)
}

type snapshotService struct {
	logger        hz_logger.Logger
	hiscoreClient *osrs.HiscoreClient
	hazelmere     *client.Hazelmere
}

func NewSnapshotService(logger hz_logger.Logger, hiscoreClient *osrs.HiscoreClient, hazelmere *client.Hazelmere) SnapshotService {
	return &snapshotService{
		logger:        logger,
		hiscoreClient: hiscoreClient,
		hazelmere:     hazelmere,
	}
}

func (ss *snapshotService) MakeSnapshotForUser(userId string) (api.HiscoreSnapshot, error) {
	user, err := ss.hazelmere.User.GetUserById(userId)
	if err != nil {
		if errors.Is(err, client.ErrUserNotFound) {
			return api.HiscoreSnapshot{}, ErrUserNotFound
		}
		return api.HiscoreSnapshot{}, errors.Join(ErrSnapshotService, err)
	}

	return ss.MakeSnapshot(user.User)
}

func (ss *snapshotService) MakeSnapshot(user api.User) (api.HiscoreSnapshot, error) {
	hiscore, err := ss.hiscoreClient.GetHiscore(user.RunescapeName)
	if err != nil {
		return api.HiscoreSnapshot{}, errors.Join(ErrSnapshotService, err)
	}

	snapshot := makeSnapshot(user, hiscore)

	snapshotRequest := api.CreateSnapshotRequest{
		Snapshot: snapshot,
	}

	response, err := ss.hazelmere.Snapshot.CreateSnapshot(snapshotRequest)
	if err != nil {
		return api.HiscoreSnapshot{}, errors.Join(ErrSnapshotService, err)
	}

	return response.Snapshot, nil
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
