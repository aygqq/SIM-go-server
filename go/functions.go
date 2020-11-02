package swagger

import (
	"fmt"
	"net/http"
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
