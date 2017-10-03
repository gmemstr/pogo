package auth

import (
	"net/http"

	"github.com/ishanjain28/pogo/common"
)

func RequireAuthorization() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		if usr := DecryptSession(r); usr != nil {
			rc.User = usr
			return nil
		}
		return &common.HTTPError{
			Message:    "Unauthorized!",
			StatusCode: http.StatusUnauthorized,
		}
	}
}

func CreateSession() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		return nil
	}
}

func DecryptSession(r *http.Request) *common.User {
	return nil
}
