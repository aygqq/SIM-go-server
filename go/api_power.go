/*
 * Power control block
 *
 * This API was created to monitor states of Power Control Block and send some commands to it.
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"../control"
)

func GetPwrBat(w http.ResponseWriter, r *http.Request) {
	var resp RespBattery
	resp.Results = 94
	resp.Status = "OK"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func GetPwrCfg(w http.ResponseWriter, r *http.Request) {
	var res RespPowercfgResults
	var resp RespPowercfg

	cfg := control.GetPowerConfig()
	res.PowerStat = cfg.PowerStat
	res.BatLevel = cfg.BatLevel
	res.Pc = cfg.Pc
	res.Wifi = cfg.Wifi
	res.Relay1 = cfg.Relay[0]
	res.Relay2 = cfg.Relay[1]
	res.Modem1 = cfg.Modem[0]
	res.Modem2 = cfg.Modem[1]

	resp.Results = &res
	resp.Status = "OK"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

// func GetPwrModemByID(w http.ResponseWriter, r *http.Request) {
// 	var resp RespPowercfg

// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusOK)

// 	str, _ := json.Marshal(resp)
// 	fmt.Fprintf(w, string(str))
// }

// func GetPwrPC(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusOK)

// 	str, _ := json.Marshal(resp)
// 	fmt.Fprintf(w, string(str))
// }

// func GetPwrRelayByID(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusOK)

// 	str, _ := json.Marshal(resp)
// 	fmt.Fprintf(w, string(str))
// }

// func GetPwrWiFi(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusOK)

// 	str, _ := json.Marshal(resp)
// 	fmt.Fprintf(w, string(str))
// }
func doSth(b bool) {

}

func SetPwrCfg(w http.ResponseWriter, r *http.Request) {
	var res RespPowercfgResults
	var resp RespPowercfg
	var state bool

	for k, v := range r.URL.Query() {
		fmt.Printf("%s: %s\n", k, v)
		if k == "pc" {
			if v[0] == "true" {
				state = true
			} else if v[0] == "false" {
				state = false
			}
			doSth(state)
		} else if k == "wifi" {
			if v[0] == "true" {
				state = true
			} else if v[0] == "false" {
				state = false
			}
			doSth(state)
		} else if k == "relay1" {
			if v[0] == "true" {
				state = true
			} else if v[0] == "false" {
				state = false
			}
			doSth(state)
		} else if k == "relay2" {
			if v[0] == "true" {
				state = true
			} else if v[0] == "false" {
				state = false
			}
			doSth(state)
		} else if k == "modem1" {
			if v[0] == "true" {
				state = true
			} else if v[0] == "false" {
				state = false
			}
			doSth(state)
		} else if k == "modem2" {
			if v[0] == "true" {
				state = true
			} else if v[0] == "false" {
				state = false
			}
			doSth(state)
		}
	}

	cfg := control.GetPowerConfig()
	res.PowerStat = cfg.PowerStat
	res.BatLevel = cfg.BatLevel
	res.Pc = cfg.Pc
	res.Wifi = cfg.Wifi
	res.Relay1 = cfg.Relay[0]
	res.Relay2 = cfg.Relay[1]
	res.Modem1 = cfg.Modem[0]
	res.Modem2 = cfg.Modem[1]

	resp.Results = &res
	resp.Status = "OK"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetPwrModemByID(w http.ResponseWriter, r *http.Request) {
	var res RespStateResults
	var resp RespState

	idx, state := parseNumberState(r)

	control.SendObjectPwr(control.OBJECT_MODEM, idx, state)

	cfg := control.GetPowerConfig()
	res.Number = idx
	res.State = cfg.Modem[idx] && state

	resp.Results = &res
	resp.Status = "OK"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetPwrPC(w http.ResponseWriter, r *http.Request) {
	var res RespStateResults
	var resp RespState

	state := parseState(r)

	control.SendObjectPwr(control.OBJECT_PC, 0, state)
	control.FlagWaitResp = true
	select {
	case read := <-control.HttpReqChan:
		if read == 1 {
			state = true
		} else {
			state = false
		}
		res.Number = 0
		res.State = state

		resp.Status = "OK"
	case <-time.After(time.Second):
		fmt.Println("No response received")

		res.Number = 0
		res.State = false

		resp.Status = "ERROR"
	}
	resp.Results = &res

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetPwrRelayByID(w http.ResponseWriter, r *http.Request) {
	var res RespStateResults
	var resp RespState

	idx, state := parseNumberState(r)

	//TODO: set this param

	cfg := control.GetPowerConfig()
	res.Number = idx
	res.State = cfg.Relay[idx] && state

	resp.Results = &res
	resp.Status = "OK"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetPwrWiFi(w http.ResponseWriter, r *http.Request) {
	var res RespStateResults
	var resp RespState

	state := parseState(r)

	//TODO: set this param

	cfg := control.GetPowerConfig()
	res.Number = 0
	res.State = cfg.Wifi && state

	resp.Results = &res
	resp.Status = "OK"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetWaitmode(w http.ResponseWriter, r *http.Request) {
	var res RespStateResults
	var resp RespState

	state := parseState(r)

	//TODO: set this param

	cfg := control.GetPowerConfig()
	res.Number = 0
	res.State = cfg.Waitmode && state

	resp.Results = &res
	resp.Status = "OK"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}
