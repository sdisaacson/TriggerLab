package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/sdisaacson/TriggerLab/internal/api"
	"github.com/sdisaacson/TriggerLab/internal/config"
	"github.com/sdisaacson/TriggerLab/internal/http/fileserver"
	"github.com/sdisaacson/TriggerLab/internal/http/handlers"
	"github.com/sdisaacson/TriggerLab/internal/http/middlewares/logreq"
	"github.com/sdisaacson/TriggerLab/internal/http/middlewares/panic"
	"github.com/sdisaacson/TriggerLab/internal/http/middlewares/webhook"
	"github.com/sdisaacson/TriggerLab/internal/metrics"
	"github.com/sdisaacson/TriggerLab/internal/pubsub"
	"github.com/sdisaacson/TriggerLab/internal/storage"
	"github.com/sdisaacson/TriggerLab/internal/version"
	"github.com/sdisaacson/TriggerLab/web"
)

const (
	readTimeout  = time.Second * 5
	writeTimeout = time.Second * 31 // IMPORTANT! Must be grater then create.maxResponseDelay value!
)

type Server struct {
	log  *zap.Logger
	echo *echo.Echo
}

func NewServer(log *zap.Logger) *Server {
	var srv = echo.New()

	srv.StdLogger = zap.NewStdLog(log)
	srv.Server.ReadTimeout = readTimeout
	srv.Server.ReadHeaderTimeout = readTimeout
	srv.Server.WriteTimeout = writeTimeout
	srv.Server.ErrorLog = srv.StdLogger
	srv.IPExtractor = NewIPExtractor()
	srv.HideBanner = true
	srv.HidePort = true

	return &Server{
		log:  log,
		echo: srv,
	}
}

func (s *Server) Register(
	ctx context.Context,
	cfg config.Config,
	rdb *redis.Client,
	stor storage.Storage,
	pub pubsub.Publisher,
	sub pubsub.Subscriber,
) error {
	registry := metrics.NewRegistry()

	s.echo.Use(
		logreq.New(s.log, []string{"/ready", "/health"}),
		panic.New(s.log),
	)

	websocketMetrics := metrics.NewWebsockets()
	if err := websocketMetrics.Register(registry); err != nil {
		return err
	}

	api.RegisterHandlers(s.echo, handlers.NewAPI(
		ctx,
		cfg,
		rdb,
		stor,
		pub,
		sub,
		registry,
		version.Version(),
		&websocketMetrics,
	))

	webhookMetrics := metrics.NewWebhooks()
	if err := webhookMetrics.Register(registry); err != nil {
		return err
	}

	var (
		wh     = webhook.New(ctx, cfg, stor, pub, &webhookMetrics)
		static = fileserver.NewHandler(http.FS(web.Content()))
	)

	s.echo.Any("/*", wh(func(c echo.Context) error { // wrap file server into webhook middleware
		if method := c.Request().Method; method == http.MethodGet || method == http.MethodHead {
			return static(c)
		}

		s.echo.HTTPErrorHandler(echo.ErrNotFound, c)

		return nil
	}))

	return nil
}

// Start the server.
func (s *Server) Start(ip string, port uint16) error {
	return s.echo.Start(ip + ":" + strconv.Itoa(int(port)))
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) error { return s.echo.Shutdown(ctx) }
