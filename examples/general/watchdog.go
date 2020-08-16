package main

import (
	"fmt"
	"github.com/coreos/go-systemd/daemon"
	log "github.com/sirupsen/logrus"
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

func StartWatchdog() error {
	timeout, err := daemon.SdWatchdogEnabled(false)
	if err != nil {
		log.Errorf("Failed to check watchdog enable state: %v", err)
		return err
	}
	fmt.Printf("Watchdog enabled: timeout=%v\n", timeout)
	if timeout > 0 {
		go watchdog(time.Second + timeout/2)
	}

	return nil
}
