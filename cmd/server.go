// main server

package cmd

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	cacheLock sync.RWMutex
	Cache     = make(map[string]CachedResponse)
)

type CachedResponse struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
	ExpiresAt  time.Time
}

func Server(port string, origin string) {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// add go routine to clear expired cache entries periodically
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			cacheLock.Lock()
			for key, cached := range Cache {
				if time.Now().After(cached.ExpiresAt) {
					delete(Cache, key)
				}
			}
			cacheLock.Unlock()
		}
	}()

	http.HandleFunc("/clearCache", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		cacheLock.Lock()
		// clear the cache
		clear(Cache)
		cacheLock.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			// proxy to origin server for non-GET requests without caching
			resp, err := client.Do(r)
			if err != nil {
				http.Error(w, "upstream request failed", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, "failed to read response body", http.StatusInternalServerError)
				return
			}
			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.WriteHeader(resp.StatusCode)
			w.Write(body)
			log.Default().Printf("Forwarded non-GET request %s %s to origin", r.Method, r.URL.Path)
			return
		}
		cacheKey := r.Method + ":" + r.URL.RequestURI()
		if cached, exists := getCachedResponse(cacheKey); exists {
			for key, values := range cached.Headers {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.Header().Set("X-Cache", "HIT")
			w.WriteHeader(cached.StatusCode)
			w.Write(cached.Body)
			log.Default().Printf("Cache HIT for %s %s. Cache memory location: %p\n", r.Method, r.URL.Path, &Cache)
		} else {
			resp, err := http.Get(origin + r.URL.Path)
			if err != nil {
				http.Error(w, "upstream request failed", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, "failed to read response body", http.StatusInternalServerError)
				return
			}
			cacheLock.Lock()
			Cache[cacheKey] = CachedResponse{
				StatusCode: resp.StatusCode,
				Headers:    resp.Header,
				Body:       body,
				ExpiresAt:  time.Now().Add(5 * time.Minute), // Cache for 5 minutes
			}
			cacheLock.Unlock()

			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.Header().Set("X-Cache", "MISS")
			w.WriteHeader(resp.StatusCode)
			w.Write(body)
			log.Default().Printf("Cache MISS for %s %s. Response cached. Cache memory location: %p\n", r.Method, r.URL.Path, &Cache)
		}
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getCachedResponse(cacheKey string) (CachedResponse, bool) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	if cached, exists := Cache[cacheKey]; exists {
		if time.Now().Before(cached.ExpiresAt) {
			return cached, true
		}
		return CachedResponse{}, false
	}
	return CachedResponse{}, false
}

