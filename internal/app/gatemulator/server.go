package gatemulator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/antnzr/gatemulator/config"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/controller"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/middleware"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/repository"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/service"
	"github.com/antnzr/gatemulator/internal/pkg/gatemulator/job"
)

type Server struct {
	config     *config.Config
	httpServer *HttpServer
	db         *sql.DB
	stopFn     sync.Once
}

func NewServer(config *config.Config, db *sql.DB) *Server {
	return &Server{config: config, db: db}
}

func (s *Server) Run(ctx context.Context) (err error) {
	var ec = make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	subscriptionRepository := repository.NewSubscriptionRepository(s.db)
	subscriberRepository := repository.NewSubscriberRepository(s.db)
	messageFileRepository := repository.NewMessageFileRepository(s.db)
	dao := repository.NewDAO(subscriptionRepository, subscriberRepository, messageFileRepository)

	jobManager := job.NewJobManager()

	subscriptionService := service.NewSubscriptionService(dao, jobManager)
	subscriberService := service.NewSubscriberService(dao)
	messageService := service.NewMessageService(dao)
	storeService := service.NewStoreService(dao)

	controller := controller.NewController(subscriptionService, subscriberService, messageService, storeService)
	authMiddleware := middleware.AuthMiddleware(subscriberService)

	middlewares := make(map[string]middleware.Middleware)
	middlewares["auth"] = authMiddleware

	s.httpServer = NewHttpServer(s.config, controller, middlewares)

	go func() {
		err := s.httpServer.Run(ctx)
		if err != nil {
			err = fmt.Errorf("HTTP server error: %w", err)
		}
		ec <- err
	}()

	// Wait for the services to exit.
	var es []string
	for i := 0; i < cap(ec); i++ {
		if err := <-ec; err != nil {
			es = append(es, err.Error())
			// If one of the services returns by a reason other than parent context canceled,
			// try to gracefully shutdown the other services to shutdown everything,
			// with the goal of replacing this service with a new healthy one.
			// NOTE: It might be a slightly better strategy to announce it as unfit for handling traffic,
			// while leaving the program running for debugging.
			if ctx.Err() == nil {
				s.Shutdown(context.Background())
			}
		}
	}
	if len(es) > 0 {
		err = errors.New(strings.Join(es, ", "))
	}
	cancel()
	return err
}

func (s *Server) Shutdown(ctx context.Context) {
	s.stopFn.Do(func() {
		s.httpServer.Shutdown(ctx)
	})
}
