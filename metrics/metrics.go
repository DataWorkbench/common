package metrics

import (
	"context"
	"net/http"

	"github.com/DataWorkbench/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Enabled bool `json:"enabled" yaml:"enabled" env:"ENABLED" validate:"-"`
	// This is used for reporting the status of server directly through
	// the HTTP address. Notice that there is a risk of leaking status
	// information if this port is exposed to the public.
	Address string `json:"address" yaml:"address" env:"ADDRESS" validate:"required_with=Enabled"`
	// HTTP URI PATH
	URLPath string `json:"url_path" yaml:"url_path" env:"URL_PATH,default=/metrics" validate:"required"`
}

// Server implements prometheus metrics server
type Server struct {
	lp  *glog.Logger
	cfg *Config
	h   *http.Server
}

// NewServer return an new Server
// NOTICE: Must set glog.Logger into the ctx by glow.WithContext
func NewServer(ctx context.Context, cfg *Config) (*Server, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	lp := glog.FromContext(ctx)
	s := &Server{
		lp:  lp,
		cfg: cfg,
	}
	return s, nil
}

func (s *Server) ListenAndServe() (err error) {
	if s == nil {
		return
	}
	mux := http.NewServeMux()
	// Expose the registered metrics via HTTP.
	mux.Handle(s.cfg.URLPath, promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
			ErrorLog:          &Logger{Output: s.lp},
		},
	))

	s.lp.Info().String("prometheus metrics server listening", s.cfg.Address+s.cfg.URLPath).Fire()

	s.h = &http.Server{Addr: s.cfg.Address, Handler: mux}

	err = s.h.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.lp.Error().Error("listen and serve prometheus metrics server error", err).Fire()
		return err
	}
	return
}

func (s *Server) Close() (err error) {
	if s == nil {
		return
	}
	s.lp.Info().Msg("waiting for prometheus metrics server close").Fire()
	if err = s.h.Close(); err != nil {
		s.lp.Error().Error("prometheus metrics server close error", err).Fire()
		return
	}
	s.lp.Info().Msg("prometheus metrics server closed").Fire()
	return
}

func (s *Server) Shutdown(ctx context.Context) (err error) {
	if s == nil {
		return
	}
	s.lp.Info().Msg("waiting for prometheus metrics server shutdown").Fire()
	if err = s.h.Shutdown(ctx); err != nil {
		s.lp.Error().Error("prometheus metrics server shutdown error", err).Fire()
		return
	}
	s.lp.Info().Msg("prometheus metrics server shutdown").Fire()
	return
}
