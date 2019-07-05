package handler

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/segmentio/ksuid"

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

	ext := h.detectExt(p.OriginalFilename)
	ct := h.detectContentType(ext)
	fn := h.genFilename(ext)

	if p.Ttl == 0 {
		presignedPostUrl, err = h.s3.CreatePresignedPostUrl(fn, ct)
	} else {
		presignedPostUrl, err = h.s3.CreatePresignedPostUrlWithTTL(fn, ct, p.Ttl)
	}

	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	} else {
		return http.StatusOK, presignedPostUrl, nil, nil
	}
}

func (h *Handler) healthcheckHandler(_ *http.Request) (int, interface{}, error) {
	ctx := h.context
	revision := ctx.Value("revision").(string)

	return http.StatusOK, entity.Health{Revision: revision}, nil
}

func (h *Handler) genFilename(ext string) string {
	uuid := ksuid.New().String()
	return fmt.Sprintf("%s%s", uuid, ext)
}

func (h *Handler) detectContentType(ext string) string {
	ct := mime.TypeByExtension(ext)
	if len(ct) > 0 {
		return ct
	} else {
		return "application/octet-stream"
	}
}

func (h *Handler) detectExt(filename string) string {
	return filepath.Ext(filename)
}
