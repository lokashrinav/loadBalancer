package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL    *url.URL
	Alive  atomic.Bool
}

type LoadBalancer struct {
	Backends []*Backend
	Current  uint32
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	backends := make([]*Backend, len(servers))
	for i, server := range servers {
		parsedURL, _ := url.Parse(server)
		backends[i] = &Backend{
			URL: parsedURL,
		}
		backends[i].Alive.Store(true)
	}
	return &LoadBalancer{Backends: backends}
}

func (lb *LoadBalancer) GetNextBackend() *Backend {
	for i := 0; i < len(lb.Backends); i++ {
		idx := atomic.AddUint32(&lb.Current, 1) % uint32(len(lb.Backends))
		backend := lb.Backends[idx]
		if backend.Alive.Load() {
			return backend
		}
	}
	return nil
}

func (lb *LoadBalancer) ForwardRequest(w http.ResponseWriter, r *http.Request) {
	backend := lb.GetNextBackend()
	if backend == nil {
		http.Error(w, "No available backend servers", http.StatusServiceUnavailable)
		return
	}

	proxyURL := backend.URL.String() + r.URL.Path
	resp, err := http.Get(proxyURL)
	if err != nil {
		http.Error(w, "Failed to reach backend", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (lb *LoadBalancer) HealthCheck(interval time.Duration) {
	for {
		time.Sleep(interval)
		for _, backend := range lb.Backends {
			resp, err := http.Get(backend.URL.String() + "/health")
			if err != nil || resp.StatusCode != http.StatusOK {
				backend.Alive.Store(false)
			} else {
				backend.Alive.Store(true)
			}
		}
	}
}

func main() {
	servers := []string{
		"http://localhost:5001",
		"http://localhost:5002",
		"http://localhost:5003",
	}

	lb := NewLoadBalancer(servers)
	go lb.HealthCheck(5 * time.Second)

	http.HandleFunc("/", lb.ForwardRequest)
	fmt.Println("Load balancer started on :8080")
	http.ListenAndServe(":8080", nil)
}
