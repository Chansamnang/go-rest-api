package pkg

import (
	"net/http"
)

type (
	response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

const (
	Success string = "success"
	Created string = "created"
	Updated string = "updated"
	Deleted string = "deleted"
	Error   string = "error"
)

func generateStatusCode(code int) int {
	if code > http.StatusNetworkAuthenticationRequired {
		code = http.StatusBadRequest
	}

	return code
}
