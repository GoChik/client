package main

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/gochik/chik"
	"github.com/gochik/chik/config"
	"github.com/gochik/chik/handlers/actor"
	"github.com/gochik/chik/handlers/datetime"
	"github.com/gochik/chik/handlers/heartbeat"
	"github.com/gochik/chik/handlers/io"
	"github.com/gochik/chik/handlers/status"
	"github.com/gochik/chik/handlers/version"
	"github.com/rs/zerolog/log"
)

// Current software version
var Version = "dev"

func main() {
	config.SetConfigFileName("client.conf")
	config.AddSearchPath("/etc/chik")
	config.AddSearchPath(".")
	err := config.ParseConfig()
	if err != nil {
		log.Warn().Msgf("Failed parsing config file: %v", err)
	}

	var server string
	config.GetStruct("server", &server)
	if server == "" {
		log.Fatal().Msg("Cannot get server from config")
	}

	log.Info().Msgf("Server: %v", server)
	controller := chik.NewController()

	ctx := context.Background()
	// Creating handlers
	go controller.Start(ctx, []chik.Handler{
		status.New(),
		io.New(),
		heartbeat.New(2 * time.Minute),
		version.New(Version),
		datetime.New(),
		actor.New(),
	})

	// Listening network
	for {
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 1 * time.Minute}, "tcp", server, &tls.Config{InsecureSkipVerify: true})
		if err == nil {
			log.Debug().Msg("New connection")
			<-controller.Connect(conn)
		}
		time.Sleep(10 * time.Second)
	}
}
