package web

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/secureworks/errors"
	"net/http"
	"strconv"
)

type HealthCheck struct {
	server *http.Server
}

func NewHealthCheck(port int) *HealthCheck {
	return &HealthCheck{
		server: &http.Server{
			Addr: ":" + strconv.Itoa(port),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}),
		},
	}
}

func (hc *HealthCheck) Start(ctx context.Context) {
	if err := hc.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Ctx(ctx).Fatal().Err(err).Msg("Health checks HTTP server failed")
	}
}

func (hc *HealthCheck) Stop(ctx context.Context) {
	if err := hc.server.Shutdown(ctx); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Health checks HTTP server shutdown failed")
	}
}
