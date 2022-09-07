package main

import (
	"github.com/luktech-labs/wax-go-sdk"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luukkk/wax-gaming-aggregator/cmd/server/handlers"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

const prefix = "WAX_GAMING_AGGREGATOR"

func main() {
	logger := log.New(os.Stdout, prefix, log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(logger); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

func run(logger *log.Logger) error {
	var cfg struct {
		conf.Version
		Web struct {
			Host            string        `conf:"default:0.0.0.0:8080"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		Wax struct {
			NodeURL string `conf:"default:https://wax.dapplica.io"`
			Proxy   string
		}
	}

	if err := conf.Parse(os.Args[1:], prefix, &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			logger.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			logger.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	logger.Printf("main: Config :\n%v\n", out)

	proxyOpt := wax.EmptyHttpOption()
	if cfg.Wax.Proxy != "" {
		proxyOpt = wax.WithProxies([]string{cfg.Wax.Proxy})
	}

	waxSdk := wax.NewSdk(cfg.Wax.NodeURL, proxyOpt)

	server := http.Server{
		Addr:         cfg.Web.Host,
		Handler:      handlers.Handlers(waxSdk),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main: Server listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main: %v : Start shutdown", sig)
	}

	return nil
}
