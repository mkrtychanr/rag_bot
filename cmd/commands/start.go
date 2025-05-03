// nolint
// Package commands provides Report API CLI handler
package commands

import (
	"context"
	"errors"
	"fmt"

	// "net/http"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/goccy/go-json"
	"github.com/mkrtychanr/rag_bot/internal/app"
	"github.com/mkrtychanr/rag_bot/internal/config"
	"github.com/mkrtychanr/rag_bot/internal/logger"
	"github.com/spf13/cobra"
)

var (
	// ErrAppRuntime is a runtime error.
	ErrAppRuntime = errors.New("error while running application") // ErrAppRuntime is a runtime error.
	startCmd      = &cobra.Command{
		Use:   "start",
		Short: "Start application",
		RunE:  startCommandRun,
	}
)

func startCommandRun(c *cobra.Command, args []string) error {
	configPath, _ := c.Flags().GetString("config")

	cfg, err := config.FindConfig(configPath)
	if err != nil {
		return fmt.Errorf("error while searching configuration. %w", err)
	}

	logger.InitLogger(cfg.Logger)

	logger.GetLogger().Info().Msg("Starting new application")
	// printConfig(*cfg)
	// Create CPU profile
	cpuprofile, _ := c.Flags().GetString("cpuprofile")
	if cpuprofile != "" {
		f, err := startCPUProfiling(cpuprofile)
		if err != nil {
			return fmt.Errorf("could not start CPU profiling. %w", err)
		}

		defer func() {
			if err := f.Close(); err != nil {
				logger.GetLogger().Err(err).Msg("Could not close cpu profile file")
			}
		}()
		defer pprof.StopCPUProfile()
	}
	// Create memory profile
	memprofile, _ := c.Flags().GetString("memprofile")
	if memprofile != "" {
		f, err := startMemoryProfiling(memprofile)
		if err != nil {
			return fmt.Errorf("could not start memory profiling. %w", err)
		}

		defer func() {
			if err := f.Close(); err != nil {
				logger.GetLogger().Err(err).Msg("Could not close memory profile file")
			}
		}()
	}

	go func() {
		logger.GetLogger().Err(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Profile.Port), nil)).Msg("pprof error")
		logger.GetLogger().Info().Msg("Bye from pprof!")
	}()

	ctx := logger.GetLogger().WithContext(context.Background())

	a, err := app.NewApp(ctx, *cfg)
	if err != nil {
		return fmt.Errorf("error while creating new application. %w", err)
	}

	return a.Run(ctx)
}

func printConfig(cfg config.Config) {
	cfgData, err := json.Marshal(cfg)
	if err != nil {
		logger.GetLogger().Debug().Err(err).Msg("Could not marshal configuration")

		return
	}

	logger.GetLogger().Debug().Msg(string(cfgData))
}

func startCPUProfiling(fileName string) (*os.File, error) {
	f, err := os.Create(filepath.Clean(fileName))
	if err != nil {
		return nil, fmt.Errorf("could not create CPU profile file %w", err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		if closeErr := f.Close(); closeErr != nil {
			logger.GetLogger().Err(closeErr).Msg("Could not close cpu profile file")
		}

		return nil, fmt.Errorf("could not start CPU profile %w", err)
	}

	return f, nil
}

func startMemoryProfiling(fileName string) (*os.File, error) {
	f, err := os.Create(filepath.Clean(fileName))
	if err != nil {
		logger.GetLogger().Err(err).Msg("Could not create memory profile")
	}

	defer func() {
		if err := f.Close(); err != nil {
			logger.GetLogger().Err(err).Msg("Could not close memory profile file")
		}
	}()

	runtime.GC()

	if err := pprof.WriteHeapProfile(f); err != nil {
		if closeErr := f.Close(); closeErr != nil {
			logger.GetLogger().Err(closeErr).Msg("Could not close cpu profile file")
		}

		return nil, fmt.Errorf("could not write memory profile. %w", err)
	}

	return f, nil
}
