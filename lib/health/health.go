package health

import (
	"log"
	"net/http"
	"time"

	"github.com/heptiolabs/healthcheck"
)

// ListenAndServe opens /live and ready endpoints
func ListenAndServe(addr string) {
	health := healthcheck.NewHandler()
	// usless, will replace later
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	// also useless
	health.AddReadinessCheck(
		"upstream-dep-dns",
		healthcheck.DNSResolveCheck("google.com", 50*time.Millisecond))
	err := http.ListenAndServe(addr, health)
	if err != nil {
		log.Fatal("Can not listen on health-checks endpoint", err)
	}
}
