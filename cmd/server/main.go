package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/Gambitier/voidkitgo/internal/config"
	"github.com/Gambitier/voidkitgo/internal/server"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
}

type CommandFlags struct {
	ConfigPath string
	Env        string
}

func NewCommandFlags() *CommandFlags {
	flags := &CommandFlags{}
	flag.StringVar(&flags.ConfigPath, "config", "default.yaml", "path to config file")
	flag.StringVar(&flags.Env, "env", string(config.Development), "environment")
	flag.Parse()
	return flags
}

// recoverPanic recovers from panics in the main thread
func recoverPanic() {
	if r := recover(); r != nil {
		logger.Errorf("Recovered from panic in main thread: %v\nStack trace:\n%s", r, debug.Stack())
		os.Exit(1)
	}
}

func main() {
	// Add panic recovery for the main thread
	defer recoverPanic()

	flags := NewCommandFlags()

	// Load configuration
	cfg, err := config.LoadConfig(logger, flags.ConfigPath, flags.Env)
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Set log level
	level, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		logger.Fatalf("Failed to parse log level: %v", err)
	}
	logger.SetLevel(level)

	// Create server instance
	srv := server.NewServer(cfg, logger)

	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server
	if err := srv.Start(ctx); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}
