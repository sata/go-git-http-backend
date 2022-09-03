package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func (s *Server) GetUploadPack(respWriter http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(respWriter, "invalid method", http.StatusBadRequest)

		return
	}

	if err := s.authenticate(req.BasicAuth()); err != nil {
		http.Error(respWriter, "invalid auth", http.StatusUnauthorized)

		return
	}

	if err := validateContentType(req, transport.UploadPackServiceName); err != nil {
		http.Error(respWriter, err.Error(), http.StatusBadRequest)

		return
	}

	packReq := &packp.UploadPackRequest{
		UploadRequest: packp.UploadRequest{
			Capabilities: &capability.List{},
			Wants:        []plumbing.Hash{},
			Shallows:     []plumbing.Hash{},
			Depth:        nil,
		},
		UploadHaves: packp.UploadHaves{
			Haves: []plumbing.Hash{},
		},
	}

	err := packReq.Decode(req.Body)
	if err != nil {
		internalErr(respWriter, err)
	}

	// otherwise validate will fail
	if packReq.Capabilities == nil {
		packReq.Capabilities = capability.NewList()
	}

	if packReq.Depth == nil {
		packReq.Depth = packp.DepthCommits(0)
	}

	ctx, cancel := context.WithTimeout(req.Context(), s.SessionTimeout)
	defer cancel()

	resp, err := s.upSession.UploadPack(ctx, packReq)
	if err != nil {
		internalErr(respWriter, err)
	}

	respWriter.Header().Add("Content-Type", fmt.Sprintf("application/x-%s-result", transport.UploadPackServiceName))
	respWriter.Header().Add("Cache-Control", "no-cache")
	respWriter.WriteHeader(http.StatusOK)

	err = resp.Encode(respWriter)
	if err != nil {
		internalErr(respWriter, err)
	}
}
