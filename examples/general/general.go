package main

import (
	"github.com/coreos/go-systemd/activation"
	"github.com/coreos/go-systemd/daemon"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net"
	"net/http"
	"os"
	"time"
)

func Main(ctx *cli.Context) error {
	addr := ctx.String("listen-address")
	listeners, err := activation.Listeners()
	if err != nil || len(listeners) == 0 {
		log.Warnf("Not support socket activation, try to listen %s", addr)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			log.Errorf("Can not listen address: %v", err)
			return err
		}
		listeners = append(listeners, l)
		log.Infof("Listen %s success", addr)
	}
	svc := ctx.String("service")
	delay := time.Duration(ctx.Int("delay-notify")) * time.Second

	since := ctx.Int("failed-until")
	if since > 0 {
		duration := time.Duration(since) * time.Second
		elapsed := UptimeSince()
		log.Infof("Time elapsed since uptime %v", elapsed)
		if elapsed < duration {
			log.Fatalf("Option failed-until set %ds, please wait %d seconds", since, duration-elapsed)
		}
	}

	if delay > 0 {
		log.Infof("Service %s will delay notify after %v", svc, delay)
		time.Sleep(delay)
	}
	if ok, err := daemon.SdNotify(false, daemon.SdNotifyReady); err != nil {
		log.Fatalf("Failed to notify %v", err)
	} else {
		if !ok {
			log.Warn("Systemd notification not supported")
		} else {
			StartWatchdog()
		}
	}

	log.Infof("Service %s started, waiting for connection", svc)
	http.HandleFunc("/", Handler)
	return http.Serve(listeners[0], nil)
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		CallerPrettyfier: nil,
	})

	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "General"
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "service", Aliases: []string{"s"}},
		&cli.StringFlag{Name: "listen-address", Value: ":8888", Aliases: []string{"l"}},
		&cli.IntFlag{Name: "delay-notify", Value: 0, Aliases: []string{"delay"}},
		&cli.IntFlag{Name: "failed-until", Value: 0},
	}
	app.Action = Main
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
