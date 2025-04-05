package service_error

import (
	"github.com/ctfloyd/hazelmere-api/src/pkg/api"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_service_error"
	"net/http"
)

var Internal = hz_service_error.ServiceError{Code: api.ErrorCodeInternal, Status: http.StatusInternalServerError}
var UserNotFound = hz_service_error.ServiceError{Code: api.ErrorCodeUserNotFound, Status: http.StatusNotFound}
