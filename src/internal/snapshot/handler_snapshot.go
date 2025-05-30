package snapshot

import (
	"errors"
	"fmt"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_handler"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_logger"
	"github.com/ctfloyd/hazelmere-worker/src/internal/common/handler"
	"github.com/ctfloyd/hazelmere-worker/src/internal/service_error"
	"github.com/ctfloyd/hazelmere-worker/src/pkg/worker_api"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type SnapshotHandler struct {
	logger  hz_logger.Logger
	service SnapshotService
}

func NewSnapshotHandler(logger hz_logger.Logger, service SnapshotService) *SnapshotHandler {
	return &SnapshotHandler{logger, service}
}

func (sh *SnapshotHandler) RegisterRoutes(mux *chi.Mux, version handler.ApiVersion) {
	if version == handler.ApiVersionV1 {
		mux.Get(fmt.Sprintf("/v1/snapshot/on-demand/{userId:%s}", hz_handler.RegexUuid), sh.CreateSnapshot)
	}
}

func (sh *SnapshotHandler) CreateSnapshot(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	snapshot, err := sh.service.MakeSnapshotForUser(userId)
	if err != nil {
		if errors.Is(err, ErrRunescapeHiscoreTimeout) {
			hz_handler.Error(w, service_error.RunescapeHiscoreTimeout, "User does not exist.")
		}
		sh.logger.ErrorArgs(r.Context(), "Could not create on demand snapshot: %v", err)
		hz_handler.Error(w, service_error.Internal, "Could not create on demand snapshot.")
		return
	}

	response := worker_api.GenerateOnDemandSnapshotResponse{SnapshotId: snapshot.Id}
	hz_handler.Ok(w, response)
}
