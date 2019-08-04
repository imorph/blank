package health

import (
	"log"
	"net/http"

	"github.com/heptiolabs/healthcheck"
)

func ListenAndServe(addr string) {
	health := healthcheck.NewHandler()
	err := http.ListenAndServe(addr, health)
	if err != nil {
		log.Fatal("Can not listen on health-checks endpoint", err)
	}
}
