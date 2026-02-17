package gatemulator

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/antnzr/gatemulator/config"
	handler "github.com/antnzr/gatemulator/internal/app/gatemulator/apperror"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/controller"
	m "github.com/antnzr/gatemulator/internal/app/gatemulator/middleware"
)

type HttpServer struct {
	config      *config.Config
	http        *http.Server
	controller  *controller.Controller
	middlewares map[string]m.Middleware
}

func NewHttpServer(
	config *config.Config,
	controller *controller.Controller,
	middlewares map[string]m.Middleware,
) *HttpServer {
	return &HttpServer{
		config:      config,
		controller:  controller,
		middlewares: middlewares,
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	router := http.NewServeMux()

	router.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})

	apiRouter := http.NewServeMux()

	logger := slog.NewLogLogger(slog.Default().Handler(), slog.LevelError)

	s.http = &http.Server{
		Addr:     fmt.Sprintf(":%v", s.config.Port),
		Handler:  router,
		ErrorLog: logger,
	}

	apiRouter.HandleFunc("POST /subscriber", handler.ErrorHandlerWrapper(s.controller.Subscriber.CreateSubscriber))
	apiRouter.HandleFunc("PATCH /subscriber/{subscriberId}", handler.ErrorHandlerWrapper(s.controller.Subscriber.UpdateSubscriber))
	apiRouter.HandleFunc("DELETE /subscriber/{subscriberId}", handler.ErrorHandlerWrapper(s.controller.Subscriber.DeleteSubscriber))

	auth := s.middlewares["auth"]
	apiRouter.Handle("POST /subscription", m.ChainMiddleware(handler.ErrorHandlerWrapper(s.controller.Subscription.Create), auth))
	apiRouter.Handle("DELETE /subscription/{subscriptionId}", m.ChainMiddleware(handler.ErrorHandlerWrapper(s.controller.Subscription.Delete), auth))
	apiRouter.Handle("PATCH /subscription/enable/{subscriptionId}", m.ChainMiddleware(handler.ErrorHandlerWrapper(s.controller.Subscription.Enable), auth))
	apiRouter.Handle("POST /subscription/scan-qr", m.ChainMiddleware(handler.ErrorHandlerWrapper(s.controller.Subscription.ScanQr), auth))

	apiRouter.Handle("POST /messages/read-recent-chat-messages/{subscriptionId}", m.ChainMiddleware(handler.ErrorHandlerWrapper(s.controller.Message.NeedReadRecentChats), auth))
	apiRouter.Handle("POST /messages/{subscriptionId}", m.ChainMiddleware(handler.ErrorHandlerWrapper(s.controller.Message.PostMessage), auth))
	apiRouter.Handle("POST /messages/messenger-sends-message", m.ChainMiddleware(handler.ErrorHandlerWrapper(s.controller.Message.MessengerSendsMessage), auth))
	apiRouter.Handle("POST /messages/read", m.ChainMiddleware(handler.ErrorHandlerWrapper(s.controller.Message.ReadMessages), auth))

	router.Handle("GET /store/{sha1}", handler.ErrorHandlerWrapper(s.controller.Store.Download))
	router.Handle("/api/", http.StripPrefix("/api", apiRouter))

	slog.Info("server is running on ::" + strconv.Itoa(s.config.Port))
	if err := s.http.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) {
	slog.Info("shutting down HTTP server")
	if s.http != nil {
		if err := s.http.Shutdown(ctx); err != nil {
			slog.Info("failed graceful shutdown of HTTP server", slog.Any("err", err))
		}
	}
}
