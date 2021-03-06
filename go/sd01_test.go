package sd01

import (
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	// speed up the test
	Timeout = 2 * time.Second
	Interval = 1 * time.Second

	announcer := NewAnnouncer("Some service", 22993)
	discoverer := NewDiscoverer("Some service")

	discoverer.Start()
	announcer.Start()

	time.Sleep(2 * time.Second)

	defer discoverer.Stop()

	services := discoverer.GetServices()
	t.Logf("Services: %+v", services)

	if len(services) != 1 {
		t.Errorf("Found %d services, expected 1", len(services))
	}

	announcer.Stop()

	time.Sleep(2 * time.Second)

	services = discoverer.GetServices()
	time.Sleep(2 * time.Second)
	services = discoverer.GetServices()

	if len(services) != 0 {
		t.Errorf("Found %d services, expected 0 after timeout", len(services))
	}
}
