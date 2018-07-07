package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/naggie/sd01/go"
)

func main() {
	// Create an announcer for our service called "example" which is running on
	// port 4001.
	a := sd01.NewAnnouncer("example", 4001)
	err := a.Start()
	if err != nil {
		panic("failed to start announcer:" + err.Error())
	}
	defer a.Stop()

	fmt.Println("announcer started")

	// Wait for interrupt.
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
