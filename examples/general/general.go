package main

import (
	"fmt"
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
	if err != nil {
		log.Warnf("Not support socket activation, try-to listen %s", addr)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			log.Errorf("Can not listen address: %v", err)
			return err
		}

		listeners = append(listeners, l)
	}
	log.Infof("Listen %s success", addr)

	if ctx.Int("delay-notify") > 0 {
		time.Sleep(time.Duration(ctx.Int("delay-notify")) * time.Second)
	}
	since:=ctx.Int("failed-until")
	if since> 0 {
		if uptimeElapse := UptimeSince();  uptimeElapse < since {
			log.Fatalf("Option failed-until set, please wait %d seconds", since-uptimeElapse)
		}
	}

	if ok, err := daemon.SdNotify(false, daemon.SdNotifyReady); err != nil {
		panic(err)
	} else {
		if !ok {
			fmt.Printf("notification not supported\n")
		} else {
			StartWatchdog()
		}
	}

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
	app.Flags = [] cli.Flag{
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
