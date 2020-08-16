package main

import (
	"fmt"
	"github.com/coreos/go-systemd/activation"
	"github.com/coreos/go-systemd/daemon"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var ticker *time.Ticker

func HelloServer(w http.ResponseWriter, req *http.Request) {
	args := req.URL.Query()
	exit := args.Get("exit")

	n, err := io.WriteString(w, fmt.Sprintf("hello world! socket activated got %s\n", exit))
	if err != nil {
		fmt.Println("failed to write string:", err)
	} else {
		fmt.Printf("%d bytes send to client %s\n", n, req.RemoteAddr)
	}

	time.Sleep(100 * time.Microsecond)
	// handle exit
	switch exit {
	case "abort":
		go func() {
			panic("got abort")
		}()
	case "normal":
		fmt.Println("Server exit with 0")
		os.Exit(0)
	case "abnormal":
		os.Exit(1)
	case "timeout":
		ticker.Stop()
	}
}

func watchdog(timeout time.Duration) {
	fmt.Println("watchdog started.")
	ticker = time.NewTicker(timeout)

	for {
		select {
		case <-ticker.C:
			ok, err := daemon.SdNotify(false, daemon.SdNotifyWatchdog)
			fmt.Println("Send keep-alive ping", ok, err)
		}
	}
}

func startWatchdog() {
	timeout, err := daemon.SdWatchdogEnabled(false)g
	if err == nil {
		fmt.Printf("Watchdog enabled: timeout=%v\n", timeout)
		go watchdog(timeout / 2)
	}
}

func Main(ctx *cli.Context) error {
	listeners, err := activation.Listeners()
	if err != nil {
		log.Warn("Not support socket activation")
		l, err := net.Listen("tcp", ctx.String("listen-address"))
		if err != nil {
			log.Errorf("Can not listen address: %v", err)
			return err
		}

		listeners = append(listeners, l)
	}

	if ok, err := daemon.SdNotify(false, daemon.SdNotifyReady); err != nil {
		panic(err)
	} else {
		if !ok {
			fmt.Printf("notification not supported\n")
		} else {
			startWatchdog()
		}
	}

	http.HandleFunc("/", HelloServer)
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
	}
	app.Action = Main
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
