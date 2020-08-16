package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func Handler(w http.ResponseWriter, req *http.Request) {
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