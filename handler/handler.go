package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cnosuke/gonzo/config"
	"github.com/cnosuke/gonzo/s3handler"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"go.uber.org/zap"
)

type Handler struct {
	router   *mux.Router
	context  context.Context
	chain    alice.Chain
	s3       *s3handler.S3Handler
	corsHost string
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHandler(ctx context.Context, conf *config.Config) *Handler {
	r := mux.NewRouter()

	h := &Handler{
		context:  ctx,
		router:   r,
		corsHost: conf.CorsHost,
		s3: s3handler.NewS3Handler(
			ctx,
			conf.AwsRegion,
			conf.BucketName,
			conf.KeyPrefix,
			conf.DefaultPresignedTTL,
		),
	}

	h.setChain(
		alice.New(
			appLoggingMiddleware(zap.L()),
		),
	)

	return h
}

func (h *Handler) Routing() *mux.Router {
	h.apiRouting()

	return h.router
}

func (h *Handler) setChain(chain alice.Chain) {
	h.chain = chain
	return
}

func (h *Handler) getChain() alice.Chain {
	return h.chain
}

func appLoggingMiddleware(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tStart := time.Now()
			next.ServeHTTP(w, r)
			tEnd := time.Now()
			t := tEnd.Sub(tStart)
			l.Sugar().Infof("[%s] %s %s", r.Method, r.URL, t.String())
		})
	}
}

type _handler struct {
	h           func(*http.Request) (int, interface{}, error)
	withHeaders func(*http.Request) (int, interface{}, map[string]string, error)
}

func (h _handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	var (
		status  int
		res     interface{}
		err     error
		headers map[string]string
	)

	if h.withHeaders != nil {
		status, res, headers, err = h.withHeaders(r)
	} else {
		status, res, err = h.h(r)
	}

	if err != nil {
		zap.S().Infof("error: %s", err)
		w.WriteHeader(status)
		encoder.Encode(res)
		return
	}

	if headers != nil {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}
	w.WriteHeader(status)

	if res != nil {
		encoder.Encode(res)
	}
	return
}
