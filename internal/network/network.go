package network

import (
	"fmt"
	"net"
)

func isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false // Port is not available
	}
	_ = ln.Close()
	return true // Port is available
}
