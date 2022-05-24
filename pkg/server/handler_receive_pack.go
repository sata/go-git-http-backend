package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func (s *Server) GetReceivePack(respWriter http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(respWriter, "invalid method", http.StatusBadRequest)

		return
	}

	if err := s.authenticate(req.BasicAuth()); err != nil {
		http.Error(respWriter, "invalid auth", http.StatusUnauthorized)

		return
	}

	if err := validateContentType(req, transport.ReceivePackServiceName); err != nil {
		http.Error(respWriter, err.Error(), http.StatusBadRequest)

		return
	}

	refReq := packp.NewReferenceUpdateRequest()

	err := refReq.Capabilities.Add(capability.ReportStatus)
	if err != nil {
		internalErr(respWriter, err)
	}

	err = refReq.Decode(req.Body)
	if err != nil {
		internalErr(respWriter, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.SessionTimeout)
	defer cancel()

	resp, err := s.rpSession.ReceivePack(ctx, refReq)
	if err != nil {
		internalErr(respWriter, err)
	}

	respWriter.Header().Add("Content-Type",
		fmt.Sprintf("application/x-%s-advertisement", transport.ReceivePackServiceName))
	respWriter.Header().Add("Cache-Control", "no-cache")
	respWriter.WriteHeader(http.StatusOK)

	err = resp.Encode(respWriter)
	if err != nil {
		internalErr(respWriter, err)
	}
}
