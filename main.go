package main

import (
	"fmt"
	"gomonitor/library/configuration"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
	cfg, err := config.LoadConfigurationJSON("config.json")
	if err != nil {
		log.Fatalf("%v", err)
	}

	url, err := cfg.GetValue("url")
	if err != nil {
		log.Fatalf("URL not defined, error: %v", err)
	}

	method, err := cfg.GetValue("method")
	if err != nil {
		method = "GET"
	}

	seconds, err := cfg.GetValue("seconds")
	if err != nil {
		seconds = "0"
	}
	s, err := strconv.Atoi(seconds)
	if err != nil {
		log.Print("not valid value for seconds, using default 60")
		s = 60
	}

	repeat, err := cfg.GetValue("repeat")
	if err != nil {
		repeat = "0"
	}
	r, err := strconv.Atoi(repeat)
	if err != nil {
		log.Print("not valid value for repeat, using default 0")
		r = 0
	}

	services := make(chan service)
	var wg sync.WaitGroup
	signalChannel := make(chan os.Signal)

	go check(services, &wg)

	i := 0
	go func() {
		for {
			wg.Add(1)
			services <- service{
				url:    url,
				method: method}
			wg.Wait()
			if r > 0 {
				i++
				if i == r {
					close(services)
					signalChannel <- syscall.SIGTERM
				}
			}
			time.Sleep(time.Duration(s) * time.Second)
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
					"url: %s | Method: %s | Response time: %f | HTTP Status: %s\n",
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

	err = resp.Body.Close()
	if err != nil {
		log.Printf(
			"an error ocurred while closing the body of the request, error: %v", err)
	}

	return resp.Status, nil
}
