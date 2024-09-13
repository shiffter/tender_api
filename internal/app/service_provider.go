package app

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq" // PostgreSQL driver.
	"log"
	"tender/internal/closer"
	"tender/internal/config"
	envCfg "tender/internal/config/env"
	handler "tender/internal/handler/api/http"
	bidsHand "tender/internal/handler/api/http/bids"
	tenderHand "tender/internal/handler/api/http/tender"
	"tender/internal/repository"
	bidsRepo "tender/internal/repository/bids"
	tenderRepo "tender/internal/repository/tender"
	userRepo "tender/internal/repository/user"
	"tender/internal/service"
	bidsSvc "tender/internal/service/bids"
	tenderSvc "tender/internal/service/tender"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	httpConfig config.HTTPConfig

	pgPool      *pgxpool.Pool
	tenderRepos repository.TenderRepos
	userRepos   repository.UsersRepos
	bidRepos    repository.BidsRepos

	tenderService service.TenderService
	bidsService   service.BidsService
	userService   service.UserService

	tenderHandler handler.TenderHandler
	bidsHandler   handler.BidsHandler
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

func (s *serviceProvider) BidsRepo(ctx context.Context) repository.BidsRepos {
	if s.bidRepos == nil {
		s.bidRepos = bidsRepo.NewBidsRepos(s.PGPool(ctx))
	}

	return s.bidRepos
}

func (s *serviceProvider) TenderService(ctx context.Context) service.TenderService {
	if s.tenderService == nil {
		s.tenderService = tenderSvc.NewTenderService(s.TenderRepo(ctx), s.UserRepo(ctx))
	}

	return s.tenderService
}

func (s *serviceProvider) BidsService(ctx context.Context) service.BidsService {
	if s.bidsService == nil {
		s.bidsService = bidsSvc.NewBidsService(s.BidsRepo(ctx), s.UserRepo(ctx), s.TenderRepo(ctx))
	}

	return s.bidsService
}

func (s *serviceProvider) TenderHandler(ctx context.Context) handler.TenderHandler {
	if s.tenderHandler == nil {
		s.tenderHandler = tenderHand.NewHandler(s.TenderService(ctx))
	}

	return s.tenderHandler
}

func (s *serviceProvider) BidsHandler(ctx context.Context) handler.BidsHandler {
	if s.bidsHandler == nil {
		s.bidsHandler = bidsHand.NewBidsHandler(s.BidsService(ctx))
	}

	return s.bidsHandler
}

//func (s *serviceProvider) runMigrations(ctx context.Context) {
//	dsn := s.pgConfig.DSN()
//	db, err := sql.Open("postgres", dsn)
//
//	instance, err := postgres.WithInstance(db, &postgres.Config{})
//	if err != nil {
//		log.Fatalf("failed to create postgres driver: %v", err)
//	}
//
//	m, err := migrate.NewWithDatabaseInstance("/home/jonny/zadanie-6105/migrations", "postgres", instance)
//	if err != nil {
//		log.Fatalf("failed to create migrate instance: %v", err)
//	}
//
//	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
//		log.Fatalf("migration failed: %v", err)
//	}
//
//	log.Println("Migration done")
//}
