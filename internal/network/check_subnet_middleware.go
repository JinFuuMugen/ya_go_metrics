package network

import (
	"net"
	"net/http"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
)

func CheckValidSubnetMiddleware(subnet string) func(http.Handler) http.Handler {
	var ipnet *net.IPNet

	if subnet != "" {
		_, parsed, err := net.ParseCIDR(subnet)
		if err != nil {
			return nil
		}
		ipnet = parsed
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if ipnet == nil {
				next.ServeHTTP(w, r)
				return
			}

			reqIP, err := resolveHeaderIP(r)
			if err != nil {
				logger.Errorf("cannot resolve request X-Real-IP header: %v", err)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			if v4 := reqIP.To4(); v4 != nil {
				reqIP = v4
			}

			if !ipnet.Contains(reqIP) {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
