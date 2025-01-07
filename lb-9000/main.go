package main

import (
	"fmt"
	appconfig "lb-9000/lb-9000/internal/config"
	"lb-9000/lb-9000/internal/orchestration"
	"lb-9000/lb-9000/internal/pool"
	"lb-9000/lb-9000/internal/proxy"
	"lb-9000/lb-9000/internal/store"
	"lb-9000/lb-9000/internal/strategy"
	"log/slog"
	"os"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	appConfig, err := appconfig.Parse("lb-9000/internal/config/config.yaml")
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	logger := slog.Default()

	orchestrator, err := orchestration.NewKubernetes(logger, appConfig)
	if err != nil {
		return fmt.Errorf("creating orchestrator: %w", err)
	}

	podPool := pool.New(
		store.Get(appConfig, logger),
		strategy.FillHoles(),
		orchestrator,
		logger,
		appConfig.RefreshRate,
	)

	proxy.Start(podPool, strconv.Itoa(appConfig.Specs.ContainerPort))

	return nil
}
