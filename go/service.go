package sd01

import (
	"fmt"
	"net"
	"time"
)

// Service provides information about a remote service.
type Service struct {
	Addr     *net.UDPAddr
	Port     int
	LastSeen time.Time
}

// String returns a human friendly representation of the Service.
func (s *Service) String() string {
	return fmt.Sprintf("%s:%d", s.Addr.IP, s.Port)
}
