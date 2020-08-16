package main

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
	"time"
)

func UptimeSince() int {
	cmd := exec.Command("uptime", "--since")
	buf, err := cmd.Output()
	if err != nil {
		log.Errorf("failed to exec uptime: %v", err)
		return -1
	}
	out := strings.TrimSpace(string(buf))

	t, err := time.Parse("2006-01-02 15:04:05", out)
	if err != nil {
		log.Errorf("failed to parse time: %v", err)
		return -1
	}
	since := int(t.Sub(time.Now()) / time.Second)
	log.Infof("system up time is %s, since=%d\n", t.Format(time.RFC3339), since)

	return since
}
