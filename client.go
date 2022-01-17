package main

import (
	"context"
	"crypto/tls"
	"flag"
	"net"
	"time"

	"github.com/gochik/chik"
	"github.com/gochik/chik/config"
	"github.com/gochik/chik/handlers/actor"
	"github.com/gochik/chik/handlers/datetime"
	"github.com/gochik/chik/handlers/heartbeat"
	"github.com/gochik/chik/handlers/heating"
	"github.com/gochik/chik/handlers/io"
	"github.com/gochik/chik/handlers/snapcast"
	"github.com/gochik/chik/handlers/status"
	"github.com/gochik/chik/handlers/systemd"
	"github.com/gochik/chik/handlers/telegram"
	"github.com/gochik/chik/handlers/version"
	"github.com/rs/zerolog/log"
)

// Current software version
var Version = "dev"
var localSearchPath string
var conf *tls.Config

func init() {
	flag.StringVar(&localSearchPath, "config", "/etc/chik", "Config file path")
}

func main() {
	flag.Parse()

	config.SetConfigFileName("client.conf")
	config.AddSearchPath(localSearchPath)
	config.AddSearchPath("/etc/chik")
	err := config.ParseConfig()
	if err != nil {
		log.Warn().Msgf("Failed parsing config file: %v", err)
	}

	var server string
	config.GetStruct("connection.server", &server)
	if server == "" {
		log.Fatal().Msg("Cannot get server from config")
	}

	var token string
	config.GetStruct("connection.token", &token)
	if token == "" {
		log.Fatal().Msg("Cannot get token from config")
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
		heating.New(),
		telegram.New(),
		systemd.New(),
		snapcast.New(),
	})

	// Listening network
	for {
		conn, err := connect(ctx, token, server)
		if err == nil {
			log.Debug().Msg("New connection")
			<-controller.Connect(conn)
		} else {
			log.Err(err).Msg("Client connection failed, retrying in 10s")
		}
		time.Sleep(10 * time.Second)
	}
}

func connect(ctx context.Context, token string, server string) (conn *tls.Conn, err error) {
	if conf == nil {
		conf, err = config.TLSConfig(ctx, token)
		if err != nil {
			log.Err(err).Msgf("Cannot get TLS config")
			return
		}
	}
	return tls.DialWithDialer(&net.Dialer{Timeout: 1 * time.Minute}, "tcp", server, conf)
}
