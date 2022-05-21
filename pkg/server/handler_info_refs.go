package server

import (
	"fmt"
	"net/http"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func (s *Server) GetInfoRefs(respWriter http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(respWriter, "invalid method", http.StatusBadRequest)

		return
	}

	// see Smart Clients section in
	// https://github.com/git/git/blob/master/Documentation/technical/http-protocol.txt
	vals := req.URL.Query()
	if len(vals) != 1 {
		http.Error(respWriter, "too many query parameters", http.StatusBadRequest)

		return
	}

	name := vals.Get("service")
	if name != transport.UploadPackServiceName && name != transport.ReceivePackServiceName {
		// If the server does not recognize the requested service name, or the
		// requested service name has been disabled by the server administrator,
		// the server MUST respond with the `403 Forbidden` HTTP status code.
		// Smart Server Response
		// https://raw.githubusercontent.com/git/git/master/Documentation/technical/http-protocol.txt
		http.Error(respWriter, "invalid service name", http.StatusForbidden)

		return
	}

	advRefs, err := buildsAdvertisedRefs(s.repo)
	if err != nil {
		internalErr(respWriter, err)

		return
	}

	respWriter.Header().Add("Content-Type", fmt.Sprintf("application/x-%s-advertisement", transport.UploadPackServiceName))
	respWriter.Header().Add("Cache-Control", "no-cache")
	respWriter.WriteHeader(http.StatusOK)

	err = advRefs.Encode(respWriter)
	if err != nil {
		internalErr(respWriter, err)

		return
	}
}

func buildsAdvertisedRefs(repo *git.Repository) (*packp.AdvRefs, error) {
	// can we not use vendor/github.com/go-git/go-git/v5/plumbing/transport/server/server.go somehow?
	advRefs := packp.NewAdvRefs()

	iter, err := repo.References()
	if err != nil {
		return nil, fmt.Errorf("repo references: %w", err)
	}

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() != plumbing.HashReference {
			return nil
		}

		advRefs.References[ref.Name().String()] = ref.Hash()

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("iter foreach: %w", err)
	}

	ref, err := repo.Reference(plumbing.HEAD, true)
	if err != nil {
		return nil, fmt.Errorf("head reference: %w", err)
	}

	h := ref.Hash()
	advRefs.Head = &h

	return advRefs, nil
}
