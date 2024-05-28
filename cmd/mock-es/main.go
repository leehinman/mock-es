package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/leehinman/mock-es/pkg/api"
	"github.com/rcrowley/go-metrics"
)

var (
	addr             string
	expire           time.Time
	percentDuplicate uint
	percentTooMany   uint
	percentNonIndex  uint
	percentTooLarge  uint
	uid              uuid.UUID
	clusterUUID      string
	metricsInterval  time.Duration
)

func init() {
	flag.StringVar(&addr, "addr", ":9200", "address to listen on ip:port")
	flag.UintVar(&percentDuplicate, "dup", 0, "percent chance StatusConflict is returned for create action")
	flag.UintVar(&percentTooMany, "toomany", 0, "percent chance StatusTooManyRequests is returned for create action")
	flag.UintVar(&percentNonIndex, "nonindex", 0, "percent chance StatusNotAcceptable is returned for create action")
	flag.UintVar(&percentTooLarge, "toolarge", 0, "percent chance StatusEntityTooLarge is returned for POST method on _bulk endpoint")
	flag.StringVar(&clusterUUID, "clusteruuid", "", "Cluster UUID of Elasticsearch we are mocking")
	flag.DurationVar(&metricsInterval, "metrics", 0, "Interval to print metrics to stdout, 0 is no metrics")

	uid = uuid.New()
	expire = time.Now().Add(24 * time.Hour)
	flag.Parse()
	if (percentDuplicate + percentTooMany + percentNonIndex) > 100 {
		log.Fatalf("Total of create action percentages must not be more than 100.\nd: %d, t:%d, n:%d", percentDuplicate, percentTooMany, percentNonIndex)
	}
	if percentTooLarge > 100 {
		log.Fatalf("percentage StatusEntityTooLarge must be less than 100")
	}
}

func main() {
	mux := http.NewServeMux()

	if metricsInterval > 0 {
		go metrics.WriteJSON(metrics.DefaultRegistry, metricsInterval, os.Stdout)
	}

	mux.Handle("/",api.NewAPIHandler(uid, clusterUUID, metrics.DefaultRegistry, expire, percentDuplicate, percentTooMany, percentNonIndex, percentTooLarge))
	if err := http.ListenAndServe(addr, mux); err != nil {
		if err != http.ErrServerClosed {
			log.Fatalf("error running HTTP server: %s", err)
		}
	}
}
