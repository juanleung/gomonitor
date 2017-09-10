package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type service struct {
	url    string
	method string
}

func main() {
	url := flag.String("u", "", "url of the services to check")
	repeat := flag.Int(
		"r", 0, "number of times the monitor will check the services")
	seconds := flag.Int("s", 60, "seconds to wait between check")
	method := flag.String("m", "GET", "HTTP request method to use")
	flag.Parse()

	services := make(chan service)
	var wg sync.WaitGroup
	signalChannel := make(chan os.Signal)

	go check(services, &wg)

	i := 0
	go func() {
		for {
			wg.Add(1)
			services <- service{
				url:    *url,
				method: *method}
			wg.Wait()
			if *repeat > 0 {
				i++
				if i == *repeat {
					close(services)
					signalChannel <- syscall.SIGINT
				}
			}
			time.Sleep(time.Duration(*seconds) * time.Second)
		}
	}()

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel
	fmt.Print("\ngomonitor terminated")
}

func check(services <-chan service, wg *sync.WaitGroup) {
	for {
		select {
		case s, open := <-services:
			if !open {
				break
			}
			start := time.Now()
			statusCode, err := makeRequest(s)
			elapsed := time.Since(start).Seconds()
			if err != nil {
				fmt.Printf(
					"an error ocurred making the request to %s | error: %v\n",
					s.url, err)
			} else {
				fmt.Printf(
					"url: %s | method: %s | response time: %f | HTTP Status: %s\n",
					s.url, strings.ToUpper(s.method), elapsed, statusCode)
			}
			wg.Done()
		}
	}
}

func makeRequest(s service) (string, error) {
	var resp *http.Response
	var err error
	method := strings.ToUpper(s.method)

	if method == http.MethodGet {
		resp, err = http.Get(s.url)
	} else if method == http.MethodPost {
		resp, err = http.Post(s.url, "text/plain", strings.NewReader(""))
	}

	if err != nil {
		return "", err
	}

	resp.Body.Close()
	return resp.Status, nil
}
