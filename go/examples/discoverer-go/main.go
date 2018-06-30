package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/naggie/sd01"
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
	for {
		select {
		case <-sig:
			return

		case service := <-sd.Discoveries():
			fmt.Println("service found at", service.String())
		}
	}
}
