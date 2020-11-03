package swagger

import (
	"fmt"
	"net/http"
	"time"

	"../control"
)

func parseState(r *http.Request) bool {
	var state bool

	for k, v := range r.URL.Query() {
		fmt.Printf("%s: %s\n", k, v)
		if k == "state" {
			if v[0] == "true" {
				state = true
			} else if v[0] == "false" {
				state = false
			}
		}
	}

	return state
}

func parseNumberState(r *http.Request) (uint8, bool) {
	var idx uint8
	var state bool

	for k, v := range r.URL.Query() {
		fmt.Printf("%s: %s\n", k, v)
		if k == "number" {
			tmp := []byte(v[0])
			idx = tmp[0] - '0'
		} else if k == "state" {
			if v[0] == "true" {
				state = true
			} else if v[0] == "false" {
				state = false
			}
		}
	}

	return idx, state
}

func doSth(b bool) {

}

func waitForResponce(resp *RespState) bool {
	var state bool = false

	control.FlagWaitResp = true
	select {
	case read := <-control.HttpReqChan:
		if read == 1 {
			state = true
		} else {
			state = false
		}
		resp.Results.Number = 0
		resp.Results.State = state

		resp.Status = "OK"
	case <-time.After(time.Second):
		fmt.Println("No response received")

		resp.Results.Number = 0
		resp.Results.State = false

		resp.Status = "ERROR"
	}
	return state
}
