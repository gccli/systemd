package main

import (
	"fmt"
	"github.com/coreos/go-systemd/daemon"
	"time"
)

var ticker *time.Ticker

func watchdog(timeout time.Duration) {
	fmt.Println("Watchdog started.")
	ticker = time.NewTicker(timeout)

	for {
		select {
		case <-ticker.C:
			ok, err := daemon.SdNotify(false, daemon.SdNotifyWatchdog)
			fmt.Println("Send keep-alive ping", ok, err)
		}
	}
}

func StartWatchdog() {
	timeout, err := daemon.SdWatchdogEnabled(false)
	if err == nil {
		fmt.Printf("Watchdog enabled: timeout=%v\n", timeout)
		go watchdog(timeout / 2)
	}
}
