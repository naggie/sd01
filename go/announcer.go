package sd01

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	// Interval between announcements.
	defaultInterval = 5 * time.Second
	minInterval     = 5 * time.Second
	maxInterval     = 60 * time.Second

	// Port is the sd01 service discovery port number.
	Port = 17823

	// maxMessageLength is the maximum allowed message length for sd01 packets.
	// This is intended to keep broadcast network traffic to a minimum.
	maxMessageLength = 64
)

// Announcer implements sd01 service announcement.
type Announcer struct {
	name     string
	port     int
	wg       *sync.WaitGroup
	stop     chan struct{}
	interval time.Duration
}

// NewAnnouncer returns a new Announcer and published beacons containing the
// supplied service name and port number.
func NewAnnouncer(name string, port int) *Announcer {
	if port < 0 || port > 65535 {
		panic("port number outside of legal range")
	}
	return &Announcer{
		name:     name,
		port:     port,
		wg:       &sync.WaitGroup{},
		interval: defaultInterval,
	}
}

// Start the Announcer. Remember to call Stop when finished.
func (a *Announcer) Start() error {
	dest, err := net.ResolveUDPAddr("udp", fmt.Sprintf("255.255.255.255:%d", Port))
	if err != nil {
		return err
	}

	local, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", local)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("sd01%s%05d", a.name, a.port)
	if len(message) > maxMessageLength {
		return fmt.Errorf("message is greater than 64 byte maximum (is %d: %s)", len(message), message)
	}

	a.wg.Add(1)
	a.stop = make(chan struct{})
	go a.run(conn, dest, message)

	return nil
}

func (a *Announcer) SetInterval(interval time.Duration) {
	if interval < minInterval || interval > maxInterval {
		panic("Specified interval out of range")
	}
	a.interval = interval
}

// Stop the Announcer.
func (a *Announcer) Stop() {
	close(a.stop)
	a.wg.Wait()
}

func (a *Announcer) run(conn *net.UDPConn, dest *net.UDPAddr, message string) {
	defer a.wg.Done()
	defer conn.Close()
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()
	for {
		select {
		case <-a.stop:
			return

		case <-ticker.C:
			if _, err := conn.WriteTo([]byte(message), dest); err != nil {
				fmt.Fprintln(os.Stderr, "sd01.announcer: failed to send beacon:", err)
			}
		}
	}
}
