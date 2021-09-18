package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/malligerearun/requests-count/web"
)

type requestCount struct {
	sync.Mutex
	requestTimeStamps map[int64]int32
}

func NewRequestCount() requestCount {
	return requestCount{
		requestTimeStamps: make(map[int64]int32),
	}
}

func (rc *requestCount)requestsCount(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Unix()
	rc.recordRequest(currentTime)
	log.Printf("requests-count::GET:: new request %v %v", rc.getRequestCount(currentTime), time.Unix(currentTime, 0))
	data := struct {		
		Count int32 `json:"Number of requests in the last 60 seconds"`
	} {
		Count: rc.getRequestCount(currentTime),
	}
	
	web.Respond(data, w)
}

func (rc *requestCount) recordRequest(timestamp int64) {
	rc.Lock()
	defer rc.Unlock()
	rc.requestTimeStamps[timestamp]++
	for seconds := range rc.requestTimeStamps {
		if seconds < (timestamp - 60) {
			delete(rc.requestTimeStamps, seconds)
		}
	}
}

func (rc *requestCount) getRequestCount(timestamp int64) int32 {
	var total int32
	rc.Lock()
	defer rc.Unlock()
	for seconds, count := range rc.requestTimeStamps {
		if seconds > (timestamp - 60) {
			total += count
		} else {
			delete(rc.requestTimeStamps, seconds)
		}
	}
	return total
}