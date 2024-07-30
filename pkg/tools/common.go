package tools

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return isLocalhost(ip)
	}

	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return isLocalhost(ip)
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return isLocalhost(ip)
	}

	return ""
}

func isLocalhost(ip string) string {
	if ip == "::1" || ip == "127.0.0.1" {
		return "127.0.0.1"
	} else {
		return ip
	}
}
