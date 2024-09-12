package app

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"tender/internal/closer"
	"tender/internal/config"
	envCfg "tender/internal/config/env"
	handler "tender/internal/handler/api/http"
	tenderHand "tender/internal/handler/api/http/tender"
	"tender/internal/repository"
	tenderRepo "tender/internal/repository/tender"
	userRepo "tender/internal/repository/user"
	"tender/internal/service"
	tenderSvc "tender/internal/service/tender"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	httpConfig config.HTTPConfig

	pgPool      *pgxpool.Pool
	tenderRepos repository.TenderRepos
	userRepos   repository.UsersRepos

	tenderService service.TenderService
	userService   service.UserService

	tenderHandler handler.TenderHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		pgCfg, err := envCfg.NewPGConfig()
		if err != nil {
			log.Fatalf("failed get pg config: %s", err.Error())
		}

		s.pgConfig = pgCfg
	}

	return s.pgConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		httpCfg, err := envCfg.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed get http config: %s", err.Error())
		}

		s.httpConfig = httpCfg
	}

	return s.httpConfig
}

func (s *serviceProvider) PGPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.Connect(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed connect to postgres: %s", err.Error())
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("failed ping db: %s", err.Error())
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		s.pgPool = pool
	}

	return s.pgPool
}

func (s *serviceProvider) TenderRepo(ctx context.Context) repository.TenderRepos {
	if s.tenderRepos == nil {
		s.tenderRepos = tenderRepo.NewTenderRepos(s.PGPool(ctx))
	}

	return s.tenderRepos
}

func (s *serviceProvider) UserRepo(ctx context.Context) repository.UsersRepos {
	if s.userRepos == nil {
		s.userRepos = userRepo.NewUserRepos(s.PGPool(ctx))
	}

	return s.userRepos
}

func (s *serviceProvider) TenderService(ctx context.Context) service.TenderService {
	if s.tenderService == nil {
		s.tenderService = tenderSvc.NewTenderService(s.TenderRepo(ctx), s.UserRepo(ctx))
	}

	return s.tenderService
}

//func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
//	if s.userService == nil {
//		s.userService = userSvc.NewUserService(s.UserRepo(ctx))
//	}
//
//	return s.userService
//}

func (s *serviceProvider) TenderHandler(ctx context.Context) handler.TenderHandler {
	if s.tenderHandler == nil {
		s.tenderHandler = tenderHand.NewHandler(s.TenderService(ctx))
	}

	return s.tenderHandler
}
