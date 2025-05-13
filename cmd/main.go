package main

import (
	requests "blcMod/internal/Requests"
	balancer "blcMod/internal/balancer"
	config "blcMod/internal/config"
	servFuncs "blcMod/internal/servers"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"
)

func main() {
	// 1. Load servers from config file
	serversConfigList := &config.ServersConfigList{Servers: make(map[string]config.ServerConfiguration)}
	serversConfigList, err := serversConfigList.GetConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Create server instances
	var servers []*servFuncs.ServerStruct
	for _, cfg := range serversConfigList.Servers {
		server := servFuncs.NewServerFunc(cfg.Name, &cfg)
		servers = append(servers, server)
		go server.Start() // Start the server in a goroutine
	}

	// 3. Create balancer
	lb := balancer.BalancerCreator(servers)
	lb.CounterWorkingServers() // Count active servers

	// 4. Simulate requests
	var wg sync.WaitGroup
	numRequests := 10
	requestInterval := 1 * time.Second
	timeout := 5 * time.Second

	fmt.Println("Starting request simulation...")
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(requestNum int) {
			defer wg.Done()
			time.Sleep(time.Duration(requestNum) * requestInterval) // Simulate delay
			startTime := time.Now()
			fmt.Printf("Request %d: Sending to balancer... ", requestNum+1)

			resp, err := requests.SendRequestToBalancer(lb, "/", timeout) // Send to root path
			if err != nil {
				log.Printf("Request %d failed: %v", requestNum+1, err)
				return
			}
			elapsedTime := time.Since(startTime)
			fmt.Printf("Response from %s (%s), Time: %s, Body: %s\n", resp.Name, resp.ServAddr, elapsedTime, resp.Body)
		}(i)
	}

	wg.Wait()
	fmt.Println("Request simulation finished.")

	// Wait a bit before shutdown
	time.Sleep(2 * time.Second)

	// Stop servers
	fmt.Println("Stopping servers...")
	for _, server := range servers {
		err := server.StopServer()
		if err != nil {
			log.Printf("Failed to stop server %s: %v", server.Name, err)
		}
	}

	fmt.Println("Program finished.")
}

func createSampleConfigFile() error {
	configList := &config.ServersConfigList{
		Servers: map[string]config.ServerConfiguration{
			"server1": {
				Name: "Server1",
				Adress: url.URL{
					Scheme: "http",
					Host:   "localhost:8080",
				},
				Readiness: true,
			},
			"server2": {
				Name: "Server2",
				Adress: url.URL{
					Scheme: "http",
					Host:   "localhost:8081",
				},
				Readiness: true,
			},
			"server3": {
				Name: "Server3",
				Adress: url.URL{
					Scheme: "http",
					Host:   "localhost:8082",
				},
				Readiness: true,
			},
		},
	}

	err := configList.WriteToConfig()
	if err != nil {
		return err
	}

	fmt.Println("Created sample config file: config.json")
	return nil
}

func init() {
	// Create a default config file if it doesn't exist
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		err := createSampleConfigFile()
		if err != nil {
			log.Fatalf("Failed to create default config file: %v", err)
		}
	}
}
