package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/malligerearun/requests-count/web"
)

type TestReq struct {
	rc requestCount
}

func NewTestReq(rc requestCount) TestReq {
	return TestReq {
		rc: rc,
	}
}

func (tr *TestReq) testRequest(w http.ResponseWriter, r *http.Request, requestsCount chan map[int64]int32) {
	currentTime := time.Now().Unix()
	rc.recordRequest(currentTime, requestsCount)
	
	log.Printf("test-request::GET:: new request %v", time.Unix(currentTime, 0))
	
	data := struct {
		Resp string `json:"Response"`
	}{
		Resp: "Response from test request",
	}

	web.Respond(data, w)
}