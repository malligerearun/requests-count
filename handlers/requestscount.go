package handlers

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/malligerearun/requests-count/web"
)

type requestCount struct {
	mu *sync.Mutex
	requestTimeStamps map[int64]int32
}

func NewRequestCount() requestCount {
	matches, err := filepath.Glob("requestCountFile")
	if err != nil || len(matches) == 0 {
		return requestCount{
			requestTimeStamps: make(map[int64]int32),
		}
	}

	r, err := os.Open("requestCountFile")
	if err != nil {
		panic(err)
	}

	decoder := gob.NewDecoder(r)
	var count map[int64]int32
	
	err = decoder.Decode(&count)
	if err != nil {
		panic(err)
	}
	
	return requestCount{
		mu: new(sync.Mutex),
		requestTimeStamps: count,
	}
}

func (rc *requestCount) requestsCount(w http.ResponseWriter, r *http.Request, requestsCount chan map[int64]int32) {
	currentTime := time.Now().Unix()
	rc.recordRequest(currentTime, requestsCount)
	reqCnt := rc.getRequestCount(currentTime)
	
	log.Printf("requests-count::GET:: new request %v %v", reqCnt, time.Unix(currentTime, 0))
	
	data := struct {
		Count int32 `json:"Number of requests in the last 60 seconds"`
	}{
		Count: reqCnt,
	}

	web.Respond(data, w)
}

func (rc *requestCount) recordRequest(timestamp int64, requestsCount chan map[int64]int32) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.requestTimeStamps[timestamp]++
	requestsCount <- rc.requestTimeStamps
	for seconds := range rc.requestTimeStamps {
		if seconds < (timestamp - 60) {
			delete(rc.requestTimeStamps, seconds)
		}
	}
}

func (rc *requestCount) getRequestCount(timestamp int64) int32 {
	var total int32
	rc.mu.Lock()
	defer rc.mu.Unlock()
	for seconds, count := range rc.requestTimeStamps {
		if seconds > (timestamp - 60) {
			total += count
		} else {
			delete(rc.requestTimeStamps, seconds)
		}
	}
	return total
}
