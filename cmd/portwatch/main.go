package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
)

const version = "0.1.0"

func main() {
	cfgPath := flag.String("config", "", "path to config file (default: built-in defaults)")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("portwatch %s\n", version)
		os.Exit(0)
	}

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	dispatcher := alert.NewDispatcher()
	handlers, err := config.DefaultHandlers(cfg)
	if err != nil {
		log.Fatalf("handler setup error: %v", err)
	}
	for _, h := range handlers {
		dispatcher.Register(h)
	}

	// Always register a log handler as fallback.
	dispatcher.Register(alert.NewLogHandler(log.Default()))

	d := daemon.New(cfg, dispatcher)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Printf("portwatch %s starting (interval: %s)", version, cfg.Interval)
	if err := d.Run(ctx); err != nil {
		log.Fatalf("daemon error: %v", err)
	}
	log.Println("portwatch stopped")
}

func loadConfig(path string) (*config.Config, error) {
	if path == "" {
		return config.Default(), nil
	}
	return config.Load(path)
}
