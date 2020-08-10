package main

import (
	"fmt"
	"github.com/coreos/go-systemd/activation"
	"github.com/coreos/go-systemd/daemon"
	"io"
	"net/http"
	"time"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	_, _ = io.WriteString(w, "hello socket activated world!\n")
}

func watchdog(timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	for {
		select {
		case <-ticker.C:
			_, _ = daemon.SdNotify(false, daemon.SdNotifyWatchdog)
		}
	}
}

func startWatchdog() {
	timeout, err := daemon.SdWatchdogEnabled(false)
	if err == nil {
		fmt.Printf("Watchdog enabled: %v\n", timeout)
		return
	}
	go watchdog(timeout / 2)
}

func main() {
	listeners, err := activation.Listeners()
	if err != nil {
		panic(err)
	}

	if len(listeners) != 1 {
		panic("Unexpected number of socket activation fds")
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
	http.Serve(listeners[0], nil)
}
