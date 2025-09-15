package main

import (
	"net/http"
	_ "net/http/pprof" // #nosec G108 - pprof is only enabled in debug mode for development
	"os"

	"log/slog"

	"github.com/0x2e/fusion/api"
	"github.com/0x2e/fusion/conf"
	"github.com/0x2e/fusion/repo"
	"github.com/0x2e/fusion/server"
	"github.com/0x2e/fusion/service/demo"
	"github.com/0x2e/fusion/service/pull"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(l)

	if conf.Debug {
		l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(l)

		go func() {
			// #nosec G114 - pprof server is only for development debugging, localhost only
			if err := http.ListenAndServe("localhost:6060", nil); err != nil {
				slog.Error("pprof server", "error", err)
				return
			}
		}()
	}

	config, err := conf.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		return
	}
	repo.Init(config.DB)

	if config.DemoMode && config.DemoModeFeeds != "" {
		seeder := demo.NewFeedSeeder(repo.NewFeed(repo.DB), repo.NewGroup(repo.DB))
		if err := seeder.SeedFeeds(config.DemoModeFeeds); err != nil {
			slog.Error("Failed to seed demo feeds", "error", err)
		}
	}

	go pull.NewPuller(repo.NewFeed(repo.DB), repo.NewItem(repo.DB), server.NewConfig(repo.NewConfig(repo.DB), config.DemoMode)).Run()

	api.Run(api.Params{
		Host:            config.Host,
		Port:            config.Port,
		PasswordHash:    config.PasswordHash,
		UseSecureCookie: config.SecureCookie,
		TLSCert:         config.TLSCert,
		TLSKey:          config.TLSKey,
		DBPath:          config.DB,
		DemoMode:        config.DemoMode,
	})
}
