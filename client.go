package main

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/gochik/chik"
	"github.com/gochik/chik/config"
	"github.com/gochik/chik/handlers/io"
	"github.com/gochik/chik/handlers/status"
	"github.com/gochik/chik/handlers/actor"
	"github.com/gochik/chik/handlers/timer"
	"github.com/gochik/chik/handlers/sunphase"
	"github.com/gochik/chik/handlers/heartbeat"
	"github.com/gochik/chik/handlers/version"
	"github.com/sirupsen/logrus"
)

// Current software version
var Version = "dev"

func main() {
	config.SetConfigFileName("client.conf")
	config.AddSearchPath("/etc/chik")
	config.AddSearchPath(".")
	err := config.ParseConfig()
	if err != nil {
		logrus.Warn("Failed parsing config file: ", err)
	}

	var server string
	config.GetStruct("server", &server)
	if server == "" {
		config.Set("server", "127.0.0.1:6767")
		config.Sync()
		logrus.Fatal("Cannot get server from config")
	}

	logrus.Debug("Server: ", server)
	controller := chik.NewController()

	// Creating handlers
	handlerList := []chik.Handler{
		status.New(),
		io.New(),
		actor.New(),
		timer.New(),
		sunphase.New(),
		heartbeat.New(2 * time.Minute),
		version.New(Version),
	}

	for _, h := range handlerList {
		controller.Start(h)
	}

	// Listening network
	for {
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 1 * time.Minute}, "tcp", server, &tls.Config{})
		if err == nil {
			logrus.Debug("New connection")
			<-controller.Connect(conn)
		}
		time.Sleep(10 * time.Second)
	}
}
