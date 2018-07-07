package sd01

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// Timeout after which a discovered service is marked offline.
	Timeout = 2 * Interval
)

// Discoverer implements sd01 service discovery and provides a list of recently
// discovered services.
type Discoverer struct {
	name       string
	services   map[string]Service
	servicesMu sync.RWMutex
	wg         sync.WaitGroup
	stop       int32
}

// NewDiscoverer returns a new Discoverer with name as the service filter.
// Matching service discoveries will be reported via the Discoveries channel.
func NewDiscoverer(name string) *Discoverer {
	return &Discoverer{
		name:     name,
		services: make(map[string]Service),
	}
}

// GetServices returns a list of recently discovered services.
func (d *Discoverer) GetServices() []Service {
	d.servicesMu.RLock()
	defer d.servicesMu.RUnlock()

	var services []Service
	now := time.Now()
	for _, s := range d.services {
		if now.Sub(s.LastSeen) < Timeout {
			services = append(services, s)
		}
	}

	return services
}

// Start the Discoverer. Remember to call Stop when finished.
func (d *Discoverer) Start() error {
	conn, err := packetConnUDP(Port)
	if err != nil {
		return err
	}

	d.wg.Add(1)
	atomic.StoreInt32(&d.stop, 0)
	go d.run(conn)

	return nil
}

// Stop the Discoverer.
func (d *Discoverer) Stop() {
	if atomic.LoadInt32(&d.stop) == 0 {
		atomic.StoreInt32(&d.stop, 1)
		d.wg.Wait()
	}
}

func (d *Discoverer) run(conn net.PacketConn) {
	defer d.wg.Done()
	defer conn.Close()

	// try to create listen socket in a loop...etc.
	buf := make([]byte, 64)

	for atomic.LoadInt32(&d.stop) == 0 {
		err := conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		if err != nil {
			fmt.Fprintln(os.Stderr, "sd01.discoverer: failed to set read deadline:", err)
			time.Sleep(Timeout)
		} else {
			buflen, addr, err := conn.ReadFrom(buf)
			if err != nil {
				if e, ok := err.(net.Error); ok && !e.Timeout() {
					fmt.Fprintln(os.Stderr, "sd01.discoverer: failed to read beacon:", err)
				}
			} else {
				if buflen == 0 || buflen > 32 {
					fmt.Fprintf(os.Stderr, "sd01.discoverer: received beacon of unsupported - length: %d, data: %s, addr: %s", buflen, string(buf[:buflen]), addr.String())
				} else if string(buf[:4]) != "sd01" {
					fmt.Fprintf(os.Stderr, "sd01.discoverer: received invalid beacon - length: %d, data: %s, addr: %s", buflen, string(buf[:buflen]), addr.String())
				} else {
					bufstr := string(buf[:buflen])
					service := bufstr[4 : len(bufstr)-5]
					portstr := bufstr[len(bufstr)-5:]
					portnum, err := strconv.Atoi(portstr)
					if err != nil {
						fmt.Fprintf(os.Stderr, "sd01.discoverer: received beacon with invalid port - length: %d, data: %s, addr: %s", buflen, string(buf[:buflen]), addr.String())
					} else if service == d.name {
						d.servicesMu.Lock()
						discovered := Service{
							Addr:     addr.(*net.UDPAddr),
							Port:     portnum,
							LastSeen: time.Now(),
						}
						d.services[discovered.String()] = discovered
						d.servicesMu.Unlock()
					}
				}
			}
		}
	}
}
