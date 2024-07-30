package middleware

import (
	"net"
	"net/http"
)

// WithCheckSubnet - check x real ip by subnet
func WithCheckSubnet(subnet *net.IPNet) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !checkSubnetAccess(subnet, r) {
				http.Error(w, "forbidden access ", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func checkSubnetAccess(subnet *net.IPNet, r *http.Request) bool {
	ipStr := r.Header.Get("X-Real-IP")
	if ipStr == "" {
		return false
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	return !subnet.Contains(ip)
}
