package service_error

import (
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_service_error"
	"github.com/ctfloyd/hazelmere-worker/src/pkg/worker_api"
	"net/http"
)

var Internal = hz_service_error.ServiceError{Code: worker_api.ErrorCodeInternal, Status: http.StatusInternalServerError}
var RunescapeHiscoreTimeout = hz_service_error.ServiceError{Code: worker_api.ErrorRunescapeHiscoreTimeout, Status: http.StatusRequestTimeout}
