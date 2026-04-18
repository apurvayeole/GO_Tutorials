//Stage 1

// package main

// import (
// 	"fmt"
// 	"net/http"
// )

// func handler(w http.ResponseWriter, r *http.Request){
// 	fmt.Fprintln(w, "Hello from go server")
// }

// func main(){
// 	http.HandleFunc("/",handler) // request st / called handler function
// 	http.ListenAndServe(":8000",nil) //starts the server
// }

package main 

import(
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	// "golang.org/x/time/rate"
	"strings"
	"sync"
	"time"
	"context"
	"github.com/redis/go-redis/v9"
)
var ctx = context.Background() //timeouts and cancellation of requests

var rdb = redis.NewClient(&redis.Options{ //creates connection
	Addr: "localhost:6379",
})

type LoadBalancer struct {
	backends []*Backend
	current  int
	mu       sync.Mutex
}
type Backend struct {
	URL     *url.URL
	Alive   bool
	mu      sync.RWMutex //allows multiple readers but one writer
}
//middleware
func authMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		
		if token == ""{
			http.Error(w,"No Token", http.StatusUnauthorized)
			return
		}

		// handle Bearer format
        if strings.HasPrefix(token, "Bearer ") {
		    token = strings.TrimPrefix(token, "Bearer ")
		}
		log.Println("TOKEN:", token)
		if token != "valid-token"{
			http.Error(w, "Invalid token" , http.StatusForbidden)
			return
		}
		next.ServeHTTP(w,r)
	})
}
// var limiter = rate.NewLimiter(1,50)
// func rateLimitMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !limiter.Allow() {
// 			http.Error(w, "Too many request", http.StatusTooManyRequests)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// // }

func chain(h http.Handler) http.Handler {
	return redisRateLimiter(authMiddleware(h)) //order matters
}
//load-balancer
func (lb *LoadBalancer) getNextBackend() *Backend {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	n := len(lb.backends)

	for i := 0; i < n; i++ { //round robin logic
		idx := (lb.current + i) % n
		if lb.backends[idx].Alive { //only healthy servers
			lb.current = idx + 1
			return lb.backends[idx]
		}
	}

	return nil // no server alive
}
func (b *Backend) checkHealth() {
	resp, err := http.Get(b.URL.String() + "/users")

	b.mu.Lock()
	defer b.mu.Unlock()

	if err != nil || resp.StatusCode != 200 {
		b.Alive = false
		log.Println(b.URL, "is DOWN")
		return
	}

	b.Alive = true
	log.Println(b.URL, "is UP")
}
func healthCheckLoop(backends []*Backend) {
	for {
		for _, b := range backends {
			go b.checkHealth()
		}
		time.Sleep(5 * time.Second)
	}
}
func redisRateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr //using IP as user id

		//increase req count
		count, err := rdb.Incr(ctx,ip).Result()
		if err != nil {
			http.Error(w,"Redis error", 500)
			return
		}
		//set expiry
		if count == 1 {
			rdb.Expire(ctx, ip, 10*time.Second)
		}

		if count > 5{
			http.Error(w,"Too many request", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w,r)
	})
}


func main(){
	// multiple instances
	url1, _ := url.Parse("http://localhost:3001")
	url2, _ := url.Parse("http://localhost:3003")

	backend1 := &Backend{URL: url1, Alive: true}
	backend2 := &Backend{URL: url2, Alive: true}

	lb := &LoadBalancer{
		backends: []*Backend{backend1, backend2},
	}

	go healthCheckLoop(lb.backends)

	http.Handle("/users", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		backend := lb.getNextBackend()

		if backend == nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(backend.URL)
		proxy.ServeHTTP(w, r)

	})))

	log.Println("Gateway with Load Balancing on 3000")
	http.ListenAndServe(":3000", nil)
}