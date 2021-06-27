package main

import (
	"context"
	"fmt"

	"github.com/seb7887/janus/internal/config"
	"github.com/seb7887/janus/internal/consumer"
	"github.com/seb7887/janus/internal/server"
	"github.com/seb7887/janus/internal/st"
	"github.com/seb7887/janus/internal/tm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	var (
		httpPort = config.GetConfig().HealthPort
		httpAddr = fmt.Sprintf(":%d", httpPort)
	)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// Healthcheck service
	g.Go(func() error {
		httpSrv := server.New(httpAddr)
		log.Infof("HTTP server running at %s", httpAddr)
		return httpSrv.Serve()
	})

	// Consumer service
	g.Go(func() error {
		return consumer.InitConsumer()
	})

	// State service
	g.Go(func() error {
		return st.StartStateListener()
	})

	// Telemetry service
	g.Go(func() error {
		return tm.StartTelemetryListener()
	})

	log.Fatal(g.Wait())
}
