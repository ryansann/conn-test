package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	serviceURLEnvVar = "SERVICE_URL"
	debugEnvVar      = "DEBUG"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("received signal %v\n", sig)
		done <- true
	}()

	url := os.Getenv(serviceURLEnvVar)
	if url == "" {
		fmt.Printf("empty %s, exiting\n", serviceURLEnvVar)
		os.Exit(1)
	}

	debug := os.Getenv(debugEnvVar)

	for {
		select {
		case <-done:
			fmt.Println("received done signal, terminating process")
			return
		case <-time.After(3 * time.Second):
			res, err := http.Get(url)
			if err != nil {
				fmt.Printf("error calling %s: %v\n", url, err)
				continue
			}

			if res.StatusCode != http.StatusOK {
				fmt.Printf("got non 200 response: %v\n", res.Status)
				continue
			} else {
				fmt.Printf("got succesful response: %v\n", res.Status)
			}

			bs, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("could not read response body: %v\n", err)
				continue
			}

			res.Body.Close()

			if debug != "" {
				fmt.Printf("got response: %s\n", string(bs))
			}
		}
	}
}
