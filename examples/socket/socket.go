package main

import (
	"fmt"
	"github.com/coreos/go-systemd/activation"
	"github.com/coreos/go-systemd/daemon"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var ticker *time.Ticker

const restartFilename = "/tmp/gosocket.restart"

func writeRestart(restart string) {
	_ = ioutil.WriteFile(restartFilename, []byte(restart), os.ModePerm)
}

func readRestart() int {
	buf, err := ioutil.ReadFile(restartFilename)
	if err != nil {
		return -1
	}

	cnt, err := strconv.Atoi(string(buf))
	if err != nil {
		return -1
	}

	return cnt
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	args := req.URL.Query()

	exit := args.Get("exit")
	restart := args.Get("restart")
	if restart != "" {
		writeRestart(restart)
	}
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
	timeout, err := daemon.SdWatchdogEnabled(false)
	if err == nil {
		fmt.Printf("Watchdog enabled: timeout=%v\n", timeout)
		go watchdog(timeout / 2)
	}
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

	cnt := readRestart()
	if cnt >= 0 {
		fmt.Printf("Restart %d\n", cnt)
		writeRestart(fmt.Sprintf("%d", cnt+1))
		os.Exit(1)
	}

	http.HandleFunc("/", HelloServer)
	http.Serve(listeners[0], nil)
}
