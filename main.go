package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/peterbourgon/ff"
	"github.com/skibish/docker-demo-app/api"
	"github.com/skibish/docker-demo-app/hits"
	"github.com/skibish/docker-demo-app/log"
)

func main() {
	l := log.NewLogger()

	fs := flag.NewFlagSet("app", flag.ExitOnError)
	var (
		port         = fs.Int("port", 80, "HTTP port")
		hostname     = fs.String("chname", "", "Set custom hostname")
		redisEnabled = fs.Bool("redis-enabled", false, "Enable Redis storage backend")
		redisURL     = fs.String("redis-url", "localhost:6379", "Redis URL")
		timeout      = fs.Duration("stop-time", 5*time.Second, "Stop time")
	)
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		l.Fatalf("Failed to get configuration: %v", err)
	}

	if *hostname == "" {
		h, err := os.Hostname()
		hostname = &h
		if err != nil {
			l.Fatalf("Failed to get hostname: %v", err)
		}
	}

	var hit *hits.Hits
	if *redisEnabled {
		hit = hits.NewHitsRedis(*redisURL)
	} else {
		hit = hits.NewHits()
	}

	a := api.New(l, *hostname, hit)

	// shutdown gracefully
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		<-sigs
		l.Print("Performing graceful shutdown")
		if *redisEnabled {
			if err := hit.CloseRedis(); err != nil {
				l.Errorf("Failed to close Redis: %v", err)
			}
		}
		select {
		case <-time.After(*timeout):
			l.Print("Closed all things")
		}

		if err := a.Shutdown(); err != nil {
			l.Errorf("Failed to shutdown server: %v", err)
		}
	}()

	l.Printf("Application is ready to listen on port: %v", *port)
	if err := a.Start(*port); err != http.ErrServerClosed {
		l.Fatalf("Server failed: %v", err)
	}

	l.Print("Exiting")
}
