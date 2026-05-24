// Package ping implements a non-blocking TCP ping function.
package ping

import (
	"net"
	"strings"
	"time"
)

// Result holds the raw outcome of a connection handshake.
type Result struct {
	Online  bool
	Latency float64
}

// Ping performs a quick TCP handshake against a target and port.
func Ping(target, port string, timeout time.Duration) Result {
	if port == "" {
		port = "22"
	}
	var address string
	if strings.Contains(target, ":") && !strings.HasPrefix(target, "[") {
		address = "[" + target + "]:" + port
	} else {
		address = net.JoinHostPort(target, port)
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return Result{Online: false}
	}
	_ = conn.Close()

	return Result{
		Online:  true,
		Latency: float64(time.Since(start).Microseconds()) / 1000.0,
	}
}
