package worker_client

import (
	"api/src/pkg/worker_api"
	"errors"
	"fmt"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_client"
)

var ErrUserNotFound = errors.Join(ErrHazelmereWorkerClient, errors.New("user not found"))

type Snapshot struct {
	prefix string
	client *hz_client.HttpClient
}

func newSnapshot(client *hz_client.HttpClient) *Snapshot {
	mappings := map[string]error{
		worker_api.ErrorCodeUserNotFound: ErrUserNotFound,
	}
	client.AddErrorMappings(mappings)

	return &Snapshot{
		prefix: "snapshot",
		client: client,
	}
}

func (ss *Snapshot) GenerateSnapshotOnDemand(userId string) (worker_api.GenerateOnDemandSnapshotResponse, error) {
	url := fmt.Sprintf("%s/on-demand/%s", ss.getBaseUrl(), userId)
	var response worker_api.GenerateOnDemandSnapshotResponse
	err := ss.client.Get(url, &response)
	if err != nil {
		return worker_api.GenerateOnDemandSnapshotResponse{}, err
	}
	return response, nil
}

func (ss *Snapshot) getBaseUrl() string {
	return fmt.Sprintf("%s/%s", ss.client.GetV1Url(), ss.prefix)
}
