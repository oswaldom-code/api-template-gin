package infrastructure

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oswaldom-code/api-template-gin/src/adapters/http/rest/handlers"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewRestHandler()
	RegisterHandlersWithOptions(router, handler, GinServerOptions{BaseURL: "/"})
	return router
}

func TestLoadPing(t *testing.T) {
	router := setupTestRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	const (
		totalRequests = 10000
		concurrency   = 100
	)

	latencies := make([]time.Duration, totalRequests)
	var (
		successCount int64
		failCount    int64
	)

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: concurrency,
			MaxConnsPerHost:     concurrency,
		},
	}

	start := time.Now()

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()

			reqStart := time.Now()
			resp, err := client.Get(server.URL + "/ping")
			latencies[idx] = time.Since(reqStart)

			if err != nil {
				atomic.AddInt64(&failCount, 1)
				return
			}
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				atomic.AddInt64(&successCount, 1)
			} else {
				atomic.AddInt64(&failCount, 1)
			}
		}(i)
	}

	wg.Wait()
	totalDuration := time.Since(start)

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	p50 := latencies[int(float64(totalRequests)*0.50)]
	p95 := latencies[int(float64(totalRequests)*0.95)]
	p99 := latencies[int(float64(totalRequests)*0.99)]

	var total time.Duration
	for _, l := range latencies {
		total += l
	}
	avg := total / time.Duration(totalRequests)

	throughput := float64(totalRequests) / totalDuration.Seconds()

	t.Logf("\n"+
		"╔══════════════════════════════════════════╗\n"+
		"║         LOAD TEST RESULTS                ║\n"+
		"╠══════════════════════════════════════════╣\n"+
		"║  Total Requests:  %6d                 ║\n"+
		"║  Concurrency:     %6d                 ║\n"+
		"║  Duration:        %10s             ║\n"+
		"║  Throughput:      %8.0f req/s         ║\n"+
		"╠══════════════════════════════════════════╣\n"+
		"║  Success:         %6d                 ║\n"+
		"║  Failures:        %6d                 ║\n"+
		"╠══════════════════════════════════════════╣\n"+
		"║  Latency avg:     %10s             ║\n"+
		"║  Latency p50:     %10s             ║\n"+
		"║  Latency p95:     %10s             ║\n"+
		"║  Latency p99:     %10s             ║\n"+
		"╚══════════════════════════════════════════╝",
		totalRequests, concurrency, totalDuration.Round(time.Millisecond),
		throughput,
		successCount, failCount,
		avg.Round(time.Microsecond), p50.Round(time.Microsecond),
		p95.Round(time.Microsecond), p99.Round(time.Microsecond),
	)

	assert.Equal(t, int64(0), failCount, "should have zero failures")
	assert.Equal(t, int64(totalRequests), successCount, "all requests should succeed")
	assert.Less(t, p99, 100*time.Millisecond, "p99 latency should be under 100ms")
}

func BenchmarkPingHandler(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewRestHandler()
	RegisterHandlersWithOptions(router, handler, GinServerOptions{BaseURL: "/"})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			router.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				b.Fatalf("unexpected status: %d", w.Code)
			}
		}
	})
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "req/s")
}

func TestLoadPingWithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewRestHandler()

	protected := router.Group("/")
	protected.Use(basicAuthorizationMiddleware)
	protected.GET("/secure-ping", func(c *gin.Context) {
		handler.Ping(c)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	const (
		totalRequests = 5000
		concurrency   = 50
	)

	var (
		rejectedCount int64
		wg            sync.WaitGroup
	)

	sem := make(chan struct{}, concurrency)
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: concurrency,
		},
	}

	start := time.Now()

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			resp, err := client.Get(server.URL + "/secure-ping")
			if err != nil {
				return
			}
			resp.Body.Close()
			if resp.StatusCode == http.StatusUnauthorized {
				atomic.AddInt64(&rejectedCount, 1)
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)
	throughput := float64(totalRequests) / duration.Seconds()

	t.Logf("\n"+
		"╔══════════════════════════════════════════╗\n"+
		"║    AUTH MIDDLEWARE LOAD TEST              ║\n"+
		"╠══════════════════════════════════════════╣\n"+
		"║  Requests:    %6d                     ║\n"+
		"║  Concurrency: %6d                     ║\n"+
		"║  Duration:    %10s                 ║\n"+
		"║  Throughput:  %8.0f req/s             ║\n"+
		"║  Rejected:    %6d (expected: all)     ║\n"+
		"╚══════════════════════════════════════════╝",
		totalRequests, concurrency,
		duration.Round(time.Millisecond), throughput,
		rejectedCount,
	)

	assert.Equal(t, int64(totalRequests), rejectedCount,
		fmt.Sprintf("all %d requests without token should be rejected with 401", totalRequests))
}
