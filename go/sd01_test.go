package sd01

import (
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	// speed up the test
	Timeout = 9

	// listen locally only so firewalls don't get in the way
	ListenAddr = [4]byte{127,0,0,1}

	announcer := NewAnnouncer("Some service", 22993)
	discoverer := NewDiscoverer("Some service")

	discoverer.Start()
	announcer.Start()

	time.Sleep(6 * time.Second)

	defer discoverer.Stop()

	services := discoverer.GetServices()
	t.Logf("Services: %+v", services)

	if len(services) != 1 {
		t.Errorf("Found %d services, expected 1", len(services))
	}

	announcer.Stop()

	time.Sleep(Timeout + time.Second)

	services = discoverer.GetServices()
	time.Sleep(Timeout + time.Second)
	services = discoverer.GetServices()

	if len(services) != 0 {
		t.Errorf("Found %d services, expected 0 after timeout", len(services))
	}
}
