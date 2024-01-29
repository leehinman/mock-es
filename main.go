package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	addr             string
	expire           time.Time
	percentDuplicate uint
	percentTooMany   uint
	percentNonIndex  uint
	percentTooLarge  uint
	uid              uuid.UUID
)

func init() {
	flag.StringVar(&addr, "a", ":9200", "address to listen on ip:port")
	flag.UintVar(&percentDuplicate, "d", 0, "percent chance StatusConflict is returned for create action")
	flag.UintVar(&percentTooMany, "t", 0, "percent chance StatusTooManyRequests is returned for create action")
	flag.UintVar(&percentNonIndex, "n", 0, "percent chance StatusNotAcceptable is returned for create action")
	flag.UintVar(&percentTooLarge, "l", 0, "percent chance StatusEntityTooLarge is returned for POST method on _bulk endpoint")
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
	mux.Handle("/", NewAPIHandler(uid, expire, percentDuplicate, percentTooMany, percentNonIndex, percentTooLarge))
	http.ListenAndServe(addr, mux)
}
