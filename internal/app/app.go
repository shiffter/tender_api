package app

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tender/internal/closer"
	"tender/internal/config"
	"time"
)

type App struct {
	serviceProvider *serviceProvider

	server *http.Server
}

func NewApp(ctx context.Context) (*App, error) {

	app := &App{}
	err := app.initDependencies(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run() error {
	defer func() {
		closer.ClosedAll()
		closer.Wait()
	}()

	return a.runHTTP()
}

func (a *App) initDependencies(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, function := range inits {
		err := function(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	a.serviceProvider = newServiceProvider()

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	router := mux.NewRouter()

	router.HandleFunc("/api/ping", a.serviceProvider.TenderHandler(ctx).Ping).Methods(http.MethodGet)
	router.HandleFunc("/api/tenders", a.serviceProvider.TenderHandler(ctx).List).Methods(http.MethodGet)
	router.HandleFunc("/api/tenders/new", a.serviceProvider.TenderHandler(ctx).Create).Methods(http.MethodPost)
	router.HandleFunc("/api/tenders/{tenderId}/status", a.serviceProvider.TenderHandler(ctx).Status).
		Methods(http.MethodGet)
	router.HandleFunc("/api/tenders/{tenderId}/status", a.serviceProvider.TenderHandler(ctx).EditStatus).
		Methods(http.MethodPut)
	router.HandleFunc("/api/tenders/my", a.serviceProvider.TenderHandler(ctx).ListForUser).Methods(http.MethodGet)

	a.server = &http.Server{
		Addr:         a.serviceProvider.HTTPConfig().Address(),
		Handler:      router,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}

	return nil
}

func (a *App) runHTTP() error {

	log.Println("starting http server on", a.serviceProvider.HTTPConfig().Address())

	if err := a.server.ListenAndServe(); err != nil {
		log.Fatalf("http server listen err: %s", err.Error())
	}

	return nil
}

func (a *App) initConfig(ctx context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}
