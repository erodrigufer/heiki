package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type Application struct {
	srv *http.Server
	// ErrorLog logs server errors.
	ErrorLog *slog.Logger
	// InfoLog informative server logger.
	InfoLog *slog.Logger
	// config centrally manages env. variables.
	config *viper.Viper
	// pgDB connection pool for postgres DB.
	pgDB *pgxpool.Pool
}

func NewApplication(ctx context.Context, requiredEnvVariables []string) (*Application, error) {
	app := new(Application)

	app.config = viper.New()
	err := app.FetchConfigValues(requiredEnvVariables)
	if err != nil {
		return nil, err
	}

	err = app.setupLoggers()
	if err != nil {
		return nil, fmt.Errorf("failed to setup the loggers: %w", err)
	}

	// err = app.setupDBConnectionPool(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to establish a DB connection pool: %w", err)
	// }

	err = app.setupServerParameters()
	if err != nil {
		return nil, fmt.Errorf("unable to setup the server's parameters: %w", err)
	}

	return app, nil
}

// StartServerWithGracefulShutdown starts a server and gracefully handles shutdowns.
// If the server receives an os.Interrupt signal the backend knows that it should
// start the process of gracefully shutting down, i.e. closing DB connections and
// closing client connections.
func (app *Application) StartServerWithGracefulShutdown(ctx context.Context) {
	go func() {
		app.InfoLog.Info("Starting serenitynow server", slog.String("server_address", app.srv.Addr))

		// Run server.
		if err := app.srv.ListenAndServe(); err != nil {
			// Error returned when server is closed, not actually an error, log to
			// info log.
			if err == http.ErrServerClosed {
				app.InfoLog.Info(err.Error())
				// An actual error happened, log to error log.
			} else {
				app.ErrorLog.Error("an error happened while executing LinstenAndServe()", slog.String("error_message", err.Error()))
			}
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// When ctx passed from main function gets cancelled with os.Interrupt signal
		// (ctx.Done() returns), this goroutine performs a shutdown.
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 15*time.Second)
		defer cancel()
		// Received an interrupt signal, shutdown.
		if err := app.srv.Shutdown(shutdownCtx); err != nil {
			// Error from closing listeners, or context timeout.
			app.ErrorLog.Error("server is not shutting down", slog.String("error_message", err.Error()))
			// An error happened while gracefully shutting down, close abruptly.
			app.srv.Close()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		app.InfoLog.Info("closing pg db connection pool")
		app.pgDB.Close()
	}()

	// Wait on all goroutines performing asynchronous shutdowns before returning.
	wg.Wait()
}

func (app *Application) setupLoggers() error {
	app.InfoLog = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	app.ErrorLog = slog.New(slog.NewJSONHandler(os.Stderr, nil))

	environmentValue, err := app.GetConfigValueString("ENVIRONMENT")
	if err != nil {
		return fmt.Errorf("unable to get config value: %w", err)
	}
	app.InfoLog = app.InfoLog.With("environment", environmentValue)
	app.ErrorLog = app.ErrorLog.With("environment", environmentValue)

	return nil
}

func (app *Application) setupDBConnectionPool(ctx context.Context) error {
	pgConnString, err := app.GetConfigValueString("PG_DB_CONN_STRING")
	if err != nil {
		return fmt.Errorf("unable to get config value: %w", err)
	}
	app.pgDB, err = pgxpool.New(ctx, pgConnString)
	if err != nil {
		return fmt.Errorf("could not establish database connection pool: %w", err)
	}
	err = app.pgDB.Ping(ctx)
	if err != nil {
		return fmt.Errorf("db unreachable, pinging db failed: %w", err)
	}

	app.InfoLog.Info("Successfully pinged db.")
	return nil
}

func (app *Application) setupServerParameters() error {
	// http.Server can only handle loggers from the old log package.
	compatibleLogger := slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, nil), slog.LevelError)

	portValue, err := app.GetConfigValueString("PORT")
	if err != nil {
		return fmt.Errorf("unable to get config value: %w", err)
	}
	app.srv = &http.Server{
		Addr:     portValue,
		ErrorLog: compatibleLogger,
		Handler:  app.routes(),
		// Time after which inactive keep-alive connections will be closed.
		IdleTimeout: time.Minute,
		// Max. time to read the header and body of a request in the server.
		ReadTimeout: 30 * time.Second,
		// Close connection if data is still being written after this time since
		// accepting the connection.
		WriteTimeout: 30 * time.Second,
	}

	return nil
}
