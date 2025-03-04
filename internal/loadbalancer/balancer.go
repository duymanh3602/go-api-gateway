package loadbalancer

import (
	"sync"
)

// Server đại diện cho một backend server
type Server struct {
	URL         string
	Connections int
	mu          sync.Mutex
}

// LoadBalancer quản lý danh sách backend servers
type LoadBalancer struct {
	servers []*Server
}

// NewLoadBalancer khởi tạo Load Balancer
func NewLoadBalancer(servers []string) *LoadBalancer {
	lb := &LoadBalancer{}
	for _, url := range servers {
		lb.servers = append(lb.servers, &Server{URL: url})
	}
	return lb
}

// SelectServer chọn server có ít kết nối nhất
func (lb *LoadBalancer) SelectServer() *Server {
	var selected *Server
	for _, server := range lb.servers {
		server.mu.Lock()
		if selected == nil || server.Connections < selected.Connections {
			selected = server
		}
		server.mu.Unlock()
	}
	// Tăng số kết nối
	selected.mu.Lock()
	selected.Connections++
	selected.mu.Unlock()
	return selected
}

// ReleaseConnection giảm số kết nối khi request hoàn tất
func (lb *LoadBalancer) ReleaseConnection(server *Server) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.Connections--
}
