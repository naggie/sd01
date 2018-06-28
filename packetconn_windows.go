// +build windows

package sd01

import (
	"fmt"
	"net"
)

// packetConnUDP calls net.ListenPacket directly on Windows for now as the
// net.FilePacketConn function is not yet implemented for windows as of v1.10.
func packetConnUDP(port int) (net.PacketConn, error) {
	return net.ListenPacket("udp", fmt.Sprintf(":%d", port))
}

// func packetConnUDP(port int) (net.PacketConn, error) {
// 	fd, err := windows.Socket(windows.AF_INET, windows.SOCK_DGRAM, windows.IPPROTO_UDP)
// 	if err != nil {
// 		windows.Close(fd)
// 		return nil, err
// 	}
// 	windows.CloseOnExec(fd)
//
// 	defer func() {
// 		if err != nil {
// 			windows.Close(fd)
// 		}
// 	}()
//
// 	if err = windows.SetsockoptInt(fd, windows.SOL_SOCKET, windows.SO_BROADCAST, 1); err != nil {
// 		return nil, err
// 	}
//
// 	if err = windows.SetsockoptInt(fd, windows.SOL_SOCKET, windows.SO_REUSEADDR, 1); err != nil {
// 		return nil, err
// 	}
//
// 	if err = windows.Bind(fd, &windows.SockaddrInet4{Port: port}); err != nil {
// 		return nil, err
// 	}
//
// 	file := os.NewFile(uintptr(fd), "")
//
// 	// NOTE: net.FilePacketConn is not yet implemented for windows and will
// 	// return an error at runtime.
// 	conn, err := net.FilePacketConn(file)
// 	if err != nil {
// 		file.Close()
// 		return nil, err
// 	}
//
// 	if err = file.Close(); err != nil {
// 		return nil, err
// 	}
//
// 	return conn, nil
// }
