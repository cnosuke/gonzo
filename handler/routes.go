package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cnosuke/gonzo/entity"
)

type route struct {
	method  string
	path    string
	handler func(http.ResponseWriter, *http.Request) (int, interface{}, error)
}

const APIPrefix = "/api/v1"

func (h *Handler) apiRouting() {
	sub := h.router.PathPrefix(APIPrefix).Subrouter()

	sub.Methods(http.MethodPost).
		Path("/create_presigned_post_url").
		Handler(h.getChain().Then(_handler{nil, h.createPresignedPostUrlHandler}))

	// Health check
	sub.Methods(http.MethodGet).
		Path("/health").
		Handler(h.getChain().Then(_handler{h.healthcheckHandler, nil}))
}

func (h *Handler) createPresignedPostUrlHandler(req *http.Request) (int, interface{}, map[string]string, error) {
	p := entity.CreatePresignedPostUrlParams{}
	if err := (json.NewDecoder(req.Body)).Decode(&p); err != nil {
		return http.StatusBadRequest, nil, nil, err
	}

	var (
		presignedPostUrl *entity.PresignedPostUrl
		err              error
	)

	if p.Ttl == 0 {
		presignedPostUrl, err = h.s3.CreatePresignedPostUrl(p.Filename, p.ContentType, p.Prefix)
	} else {
		presignedPostUrl, err = h.s3.CreatePresignedPostUrlWithTTL(p.Filename, p.ContentType, p.Prefix, p.Ttl)
	}

	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	} else {
		headers := map[string]string{"Access-Control-Allow-Origin": h.corsHost}
		return http.StatusOK, presignedPostUrl, headers, nil
	}
}

func (h *Handler) healthcheckHandler(_ *http.Request) (int, interface{}, error) {
	ctx := h.context
	revision := ctx.Value("revision").(string)

	return http.StatusOK, entity.Health{Revision: revision}, nil
}
