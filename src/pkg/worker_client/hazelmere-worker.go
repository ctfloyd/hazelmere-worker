package worker_client

import (
	"errors"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_client"
)

var ErrHazelmereWorkerClient = errors.New("generic hazelmere worker error")

type HazelmereWorker struct {
	Snapshot *Snapshot
}

func NewHazelmereWorker(client *hz_client.HttpClient) *HazelmereWorker {
	return &HazelmereWorker{
		Snapshot: newSnapshot(client),
	}
}
