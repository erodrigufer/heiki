package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/erodrigufer/serenitynow/internal/server"
)

func main() {
	ctx := context.Background()

	requiredEnvVariables := []string{"PORT", "ENVIRONMENT", "SQLITE_PATH"}
	if err := run(ctx, requiredEnvVariables); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// run encapsulates the web application.
// It can be used to create integration/unit tests using the `go test` framework.
func run(ctx context.Context, requiredEnvVariables []string) error {
	// Cancel context if application receives a SIGNINT signal, and use cancelled
	// context to start a graceful shutdown of the application.
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	app, err := server.NewApplication(ctx, requiredEnvVariables)
	if err != nil {
		return fmt.Errorf("unable to initialize a new application: %w", err)
	}
	app.StartServerWithGracefulShutdown(ctx)

	return nil
}
