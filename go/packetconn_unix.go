// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package sd01

import (
	"net"
	"os"

	"golang.org/x/sys/unix"
)

func packetConnUDP(port int) (net.PacketConn, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, unix.IPPROTO_UDP)
	if err != nil {
		unix.Close(fd)
		return nil, err
	}
	unix.CloseOnExec(fd)

	defer func() {
		if err != nil {
			unix.Close(fd)
		}
	}()

	if err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_BROADCAST, 1); err != nil {
		return nil, err
	}

	if err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		return nil, err
	}

	if err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		return nil, err
	}

	if err = unix.Bind(fd, &unix.SockaddrInet4{Port: port, Addr: ListenAddr}); err != nil {
		return nil, err
	}

	file := os.NewFile(uintptr(fd), "")

	conn, err := net.FilePacketConn(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return conn, nil
}
