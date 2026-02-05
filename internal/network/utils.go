package network

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

var errFailedToParseIP = errors.New("cannot parse ip from http header")

// OutboundIPTo returns the local IP that would be used to reach targetAddr.
func OutboundIPTo(targetAddr string) (net.IP, error) {
	host, port, err := net.SplitHostPort(targetAddr)
	if err != nil {
		return nil, fmt.Errorf("split host:port: %w", err)
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf("lookup target host ip: %w", err)
	}

	var targetIP net.IP
	for _, ip := range ips {
		if v4 := ip.To4(); v4 != nil {
			targetIP = v4
			break
		}
	}
	if targetIP == nil && len(ips) > 0 {
		targetIP = ips[0]
	}
	if targetIP == nil {
		return nil, fmt.Errorf("no IPs found for host %q", host)
	}

	conn, err := net.Dial("udp", net.JoinHostPort(targetIP.String(), port))
	if err != nil {
		return nil, fmt.Errorf("udp dial: %w", err)
	}
	defer conn.Close()

	la, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok || la.IP == nil {
		return nil, fmt.Errorf("cannot get local udp addr")
	}

	if v4 := la.IP.To4(); v4 != nil {
		return v4, nil
	}
	return la.IP, nil
}

func resolveHeaderIP(r *http.Request) (net.IP, error) {

	ipStr := r.Header.Get("X-Real-IP")
	ip := net.ParseIP(ipStr)

	if ip == nil {
		return nil, errFailedToParseIP
	}
	return ip, nil

}
