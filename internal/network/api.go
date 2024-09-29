package network

func FindFreePort(startPort int) int {
	for port := startPort; port <= 65535; port++ {
		if isPortAvailable(port) {
			return port
		}
	}
	return -1 // No free port found
}
