package swagger

import (
	"fmt"
	"net/http"
	"time"

	"../control"
)

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

func parseNumberImei(r *http.Request) (uint8, string) {
	var idx uint8
	var str string

	for k, v := range r.URL.Query() {
		fmt.Printf("%s: %s\n", k, v)
		if k == "number" {
			tmp := []byte(v[0])
			idx = tmp[0] - '0'
		} else if k == "imei" {
			str = v[0]
		}
	}

	return idx, str
}

func parseNumberSim(r *http.Request) (uint8, uint8) {
	var idx uint8
	var num uint8

	for k, v := range r.URL.Query() {
		fmt.Printf("%s: %s\n", k, v)
		if k == "number" {
			tmp := []byte(v[0])
			idx = tmp[0] - '0'
		} else if k == "sim_num" {
			tmp := []byte(v[0])
			num = tmp[0] - '0'
		}
	}

	return idx, num
}

func parseNumberPhoneSms(r *http.Request) (uint8, string, string) {
	var idx uint8
	var phone string
	var sms string

	for k, v := range r.URL.Query() {
		fmt.Printf("%s: %s\n", k, v)
		if k == "number" {
			tmp := []byte(v[0])
			idx = tmp[0] - '0'
		} else if k == "phone" {
			phone = v[0]
		} else if k == "message" {
			sms = v[0]
		}
	}

	return idx, phone, sms
}

func waitForResponce() (string, bool) {
	var ret bool
	var status string

	control.FlagWaitResp = true
	select {
	case read := <-control.HttpReqChan:
		// if read == 1 {
		// 	status = "OK"
		// 	ret = true
		// } else {
		// 	status = "ERROR"
		// 	ret = false
		// }
		fmt.Printf("Chanel recv %d\n", read)
		status = "OK"
		ret = true
	case <-time.After(time.Second):
		fmt.Println("No response received")
		status = "ERROR"
		ret = false
	}
	return status, ret
}
