package sd01

import (
	"testing"
)

func TestDiscovery(t *testing.T) {
	announcer := NewAnnouncer("foobar", 22993)
	discoverer := NewDiscoverer("foobar")

	announcer.Start()
	discoverer.Start()

	defer announcer.Stop()
	defer discoverer.Stop()

	services := discoverer.GetServices(true)
	t.Logf("Services: %v", services)

	if len(services) != 1 {
		t.Errorf("Found %v services, expected 1", len(services))
	}
}
