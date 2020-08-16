package main

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
	"time"
)

func UptimeSince() time.Duration {
	cmd := exec.Command("uptime", "--since")
	buf, err := cmd.Output()
	if err != nil {
		log.Errorf("failed to exec uptime: %v", err)
		return 0
	}
	out := strings.TrimSpace(string(buf))

	t, err := time.ParseInLocation("2006-01-02 15:04:05", out, time.UTC)
	if err != nil {
		log.Errorf("failed to parse time: %v", err)
		return 0
	}
	elapsed := time.Now().Sub(t)
	log.Infof("System uptime is %s, elapsed=%v\n", t.Format(time.RFC3339), elapsed)

	return elapsed
}
