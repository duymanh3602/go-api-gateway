package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go-api-gateway/internal/handlers"
	"go-api-gateway/internal/loadbalancer"
	"go-api-gateway/internal/middleware"
)

func main() {
	backendServers := []string{
		"http://localhost:5001",
		"http://localhost:5002",
		"http://localhost:5003",
	}

	lb := loadbalancer.NewLoadBalancer(backendServers)
	rateLimiter := middleware.NewRateLimiter(5, 10*time.Second) // Giới hạn 5 request mỗi 10 giây

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.ProxyRequest(lb))

	handler := middleware.LoggerMiddleware(rateLimiter.Middleware(mux))

	fmt.Println("Load Balancer running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
