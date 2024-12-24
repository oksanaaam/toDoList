package loadbalancer

import (
	"log"
	"net/http"
	"sync/atomic"
)

func init() {
	log.Println("Load Balancer package loaded successfully.")
}

type LoadBalancer struct {
	servers []string
	counter uint64
}

// NewLoadBalancer create a new load balancer
func NewLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{servers: servers}
}

// getNextServer select next server
func (lb *LoadBalancer) getNextServer() string {
	// atomically increment the counter and get the next server
	idx := atomic.AddUint64(&lb.counter, 1) % uint64(len(lb.servers))
	return lb.servers[idx]
}

// ServeHTTP method for processing HTTP requests and proxying them to servers
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.getNextServer()

	// Create a channel for the server response
	ch := make(chan string)

	// Asynchronously handle the redirect in a goroutine
	go lb.handleRequest(w, r, server, ch)

	// Wait for the response
	response := <-ch

	// Respond back to the client with the result
	http.Redirect(w, r, response, http.StatusTemporaryRedirect)
}

// handleRequest asynchronously handles the proxying of a request
func (lb *LoadBalancer) handleRequest(w http.ResponseWriter, r *http.Request, server string, ch chan string) {
	log.Printf("Routing request to server: %s", server)

	// Redirect the request to the selected server
	proxyURL := "http://" + server + r.URL.Path
	ch <- proxyURL // Send the redirect URL to the channel
}
