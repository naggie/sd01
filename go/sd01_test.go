package sd01

import (
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	// speed up the test
	Timeout = 9

	announcer := NewAnnouncer("Some service", 22993)
	discoverer := NewDiscoverer("Some service")

	announcer.Start()
	discoverer.Start()

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
