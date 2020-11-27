package swagger

import (
	"fmt"
	"net/http"
	"time"

	"../control"
)

func parseNumberState(r *http.Request) (uint8, bool, uint8) {
	var idx uint8
	var state bool
	var err uint8

	for k, v := range r.URL.Query() {
		if k == "number" {
			tmp := []byte(v[0])
			idx = tmp[0] - '0'
			if idx > 2 || idx < 1 {
				err = 1
			}
		} else if k == "state" {
			if v[0] == "true" {
				state = true
			} else if v[0] == "false" {
				state = false
			} else {
				err = 1
			}
		}
	}

	return idx, state, err
}

func parseNumberImei(r *http.Request) (uint8, string, uint8) {
	var idx uint8
	var imei string
	var err uint8

	for k, v := range r.URL.Query() {
		if k == "number" {
			tmp := []byte(v[0])
			idx = tmp[0] - '0'
			if idx > 2 || idx < 1 {
				err = 1
			}
		} else if k == "imei" {
			imei = v[0]
			if len(imei) != control.IMEI_SIZE {
				err = 1
			}
		}
	}

	return idx, imei, err
}

func parseNumberSim(r *http.Request) (uint8, uint8, uint8) {
	var idx uint8
	var num uint8
	var err uint8

	for k, v := range r.URL.Query() {
		if k == "number" {
			tmp := []byte(v[0])
			idx = tmp[0] - '0'
			if idx > 2 || idx < 1 {
				err = 1
			}
		} else if k == "sim_num" {
			tmp := []byte(v[0])
			num = tmp[0] - '0'
			if num > 4 || num < 1 {
				err = 1
			}
		}
	}

	return idx, num, err
}

func parseNumberPhoneSms(r *http.Request) (uint8, string, string, uint8) {
	var idx uint8
	var phone string
	var sms string
	var err uint8

	for k, v := range r.URL.Query() {
		if k == "number" {
			tmp := []byte(v[0])
			idx = tmp[0] - '0'
			if idx > 2 || idx < 1 {
				err = 1
			}
		} else if k == "phone" {
			phone = v[0]
			if len(phone) > control.PHONE_SIZE {
				err = 1
			}
		} else if k == "message" {
			sms = v[0]
		}
	}

	return idx, phone, sms, err
}

func waitForResponce() (string, bool) {
	var ret bool
	var status string

	control.FlagHTTPWaitResp = true
	select {
	case read := <-control.HTTPReqChan:
		//! COM now in echo mode, so that "read" value doesn't matter
		// if read == 1 {
		// 	status = "OK"
		// 	ret = true
		// } else {
		// 	status = "EXECUTE_ERROR"
		// 	ret = false
		// }
		fmt.Printf("Chanel recv %d\n", read)
		status = "OK"
		ret = true
	case <-time.After(time.Second):
		fmt.Println("No response received")
		status = "EXECUTE_ERROR"
		ret = false
	}
	return status, ret
}
