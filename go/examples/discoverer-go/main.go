package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/naggie/sd01/go"
)

func main() {
	// Create a discoverer which listens for announcer packets describing the
	// "example" service.
	sd := sd01.NewDiscoverer("example")
	err := sd.Start()
	if err != nil {
		panic("failed to start discoverer:" + err.Error())
	}
	defer sd.Stop()

	fmt.Println("discoverer started")

	// Wait for interrupt or discoveries.
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-sig:
			return

		case <-ticker.C:
			services := sd.GetServices()
			fmt.Printf("found %d services:\n", len(services))
			for _, service := range services {
				fmt.Printf("  - %s\n", service.String())
			}
		}
	}
}
