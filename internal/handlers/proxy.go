package handlers

import (
	"io"
	"log"
	"net/http"

	"go-api-gateway/internal/loadbalancer"
)

// ProxyRequest chuyển tiếp request tới backend server
func ProxyRequest(lb *loadbalancer.LoadBalancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server := lb.SelectServer()
		log.Printf("Forwarding request to: %s", server.URL)

		resp, err := http.Get(server.URL)
		if err != nil {
			http.Error(w, "Backend service unavailable", http.StatusServiceUnavailable)
			lb.ReleaseConnection(server)
			return
		}
		defer resp.Body.Close()

		// Copy response từ backend về client
		body, _ := io.ReadAll(resp.Body)
		w.WriteHeader(resp.StatusCode)
		w.Write(body)

		// Giảm số kết nối sau khi hoàn thành
		lb.ReleaseConnection(server)
	}
}
