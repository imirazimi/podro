package main

import (
	"context"
	"flag"
	"interview/adapter"
	"interview/pkg"
	"interview/podroapp"
	"interview/scheduler"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	config, err := pkg.GetConfig()
	if err != nil {
		panic(err)
	}
	err = pkg.InitLogger(config.Logger)
	if err != nil {
		panic(err)
	}
}

func main() {
	pkg.Logger.Info("Starting application")
	config, err := pkg.GetConfig()
	if err != nil {
		panic(err)
	}

	postgres := adapter.NewSQLDB(config.SQLDB)
	if err := postgres.Start(); err != nil {
		panic(err)
	}
	migrateDown := flag.Bool("migrate-down", false, "perform database migration down")
	flag.Parse()
	migrator := pkg.NewMigrator(config.Migrator)
	if *migrateDown {
		if err := migrator.Down(config.SQLDB); err != nil {
			panic(err)
		}
	} else {
		if err := migrator.Up(config.SQLDB); err != nil {
			panic(err)
		}
	}

	fiber := adapter.NewHTTPServer(config.HTTPServer)
	gocron, err := adapter.NewScheduler()
	if err != nil {
		panic(err)
	}
	kavenegar := adapter.NewOTP()

	podroService := podroapp.SetupService(postgres, fiber)
	podroClient := podroapp.SetupClient(podroService, kavenegar)
	podroService.SetClient(podroClient)

	scheduler.SetupService(gocron, podroClient)

	// Error channel for collecting errors from goroutines
	errCh := make(chan error, 2)

	// Start HTTP server
	go func() {
		if err := fiber.Start(); err != nil {
			pkg.Logger.Error("HTTP server error: " + err.Error())
			errCh <- err
		}
	}()

	// Start scheduler
	go func() {
		gocron.Start()
	}()

	// Wait for shutdown signal or error
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Wait for either error or shutdown signal
	select {
	case err := <-errCh:
		pkg.Logger.Error("Service error occurred: " + err.Error())
	case <-stopCh:
		pkg.Logger.Info("Received shutdown signal")
	}

	// Cleanup channels
	signal.Stop(stopCh)
	close(stopCh)
	close(errCh)

	pkg.Logger.Info("Shutting down application")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := fiber.Shutdown(ctx); err != nil {
		pkg.Logger.Error("Error shutting down HTTP server: " + err.Error())
	} else {
		pkg.Logger.Info("HTTP server shutdown successfully")
	}

	pkg.Logger.Info("Shutting down scheduler")
	gocron.ShutDown()
	pkg.Logger.Info("Scheduler shutdown successfully")

	pkg.Logger.Info("Shutting down database connection")
	postgres.ShutDown()
	pkg.Logger.Info("Database connection closed successfully")

	pkg.Logger.Info("Application stopped gracefully")
}
