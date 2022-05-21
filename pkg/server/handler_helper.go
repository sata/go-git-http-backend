package server

import (
	"fmt"
	"net/http"
)

var ErrInvalidContentType = fmt.Errorf("invalid content type")

func validateContentType(r *http.Request, service string) error {
	if r.Header.Get("Content-Type") != fmt.Sprintf("application/x-%s-request", service) {
		return ErrInvalidContentType
	}

	return nil
}

func internalErr(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintf("internal error: %s", err), http.StatusInternalServerError)
}
