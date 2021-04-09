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
	"log"
	"net/http"

	"../control"
)

func GetPwrCfg(w http.ResponseWriter, r *http.Request) {
	var res RespPowercfgResults
	var resp RespPowercfg

	if control.USBChanWaitNotBusy(1000) {
		control.FlagHTTPWaitResp = true
		control.SendCommand(control.CMD_GET_CONFIG, true)
		status, ret := waitForResponce(1)

		if ret {
			cfg := &control.PowerSt
			res.PowerStat = cfg.PowerStat
			res.BatLevel = cfg.BatLevel
			res.Pc = cfg.Pc
			res.Wifi = cfg.Wifi
			res.Relay1 = cfg.Relay[0]
			res.Relay2 = cfg.Relay[1]
			res.Modem1 = cfg.Modem[0]
			res.Modem2 = cfg.Modem[1]
			res.SimNum1 = control.ModemSt[0].SimNum
			res.SimNum2 = control.ModemSt[1].SimNum
			resp.Results = &res
		}
		resp.Status = status
	} else {
		resp.Status = "CHANEL_BUSY"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetPwrCfg(w http.ResponseWriter, r *http.Request) {
	var res RespPowercfgResults
	var resp RespPowercfg
	var err uint8

	newCfg := control.GetConfigFile()

	if !control.USBChanWaitNotBusy(1000) {
		err = 2
	} else {
		for k, v := range r.URL.Query() {
			log.Printf("%s: %s\n", k, v)
			if k == "pc" {
				if v[0] == "true" {
					newCfg.Power.Pc = true
				} else if v[0] == "false" {
					newCfg.Power.Pc = false
				} else {
					err = 1
				}
			} else if k == "wifi" {
				if v[0] == "true" {
					newCfg.Power.Wifi = true
				} else if v[0] == "false" {
					newCfg.Power.Wifi = false
				} else {
					err = 1
				}
			} else if k == "relay1" {
				if v[0] == "true" {
					newCfg.Power.Relay[0] = true
				} else if v[0] == "false" {
					newCfg.Power.Relay[0] = false
				} else {
					err = 1
				}
			} else if k == "relay2" {
				if v[0] == "true" {
					newCfg.Power.Relay[1] = true
				} else if v[0] == "false" {
					newCfg.Power.Relay[1] = false
				} else {
					err = 1
				}
			} else if k == "modem1" {
				if v[0] == "true" {
					newCfg.Power.Modem[0] = true
				} else if v[0] == "false" {
					newCfg.Power.Modem[0] = false
				} else {
					err = 1
				}
			} else if k == "modem2" {
				if v[0] == "true" {
					newCfg.Power.Modem[1] = true
				} else if v[0] == "false" {
					newCfg.Power.Modem[1] = false
				} else {
					err = 1
				}
			} else if k == "simnum1" {
				tmp := []byte(v[0])
				newCfg.SimNum[0] = tmp[0] - '0'
				if newCfg.SimNum[0] > 4 {
					err = 1
				}
			} else if k == "simnum2" {
				tmp := []byte(v[0])
				newCfg.SimNum[1] = tmp[0] - '0'
				if newCfg.SimNum[1] > 4 {
					err = 1
				}
			}
		}
	}

	if err == 0 {
		control.FlagHTTPWaitResp = true
		control.SendConfig(newCfg)
		status, ret := waitForResponce(1)
		if ret {
			cfg := &control.PowerSt
			res.PowerStat = cfg.PowerStat
			res.BatLevel = cfg.BatLevel
			res.Pc = cfg.Pc
			res.Wifi = cfg.Wifi
			res.Relay1 = cfg.Relay[0]
			res.Relay2 = cfg.Relay[1]
			res.Modem1 = cfg.Modem[0]
			res.Modem2 = cfg.Modem[1]
			res.SimNum1 = control.ModemSt[0].SimNum
			res.SimNum1 = control.ModemSt[1].SimNum

			resp.Results = &res
		}
		resp.Status = status
	} else if err == 1 {
		resp.Status = "INVALID_REQUEST"
	} else if err == 2 {
		resp.Status = "CHANEL_BUSY"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetPwrModemByID(w http.ResponseWriter, r *http.Request) {
	var res RespStateResults
	var resp RespState

	if control.USBChanWaitNotBusy(1000) {
		idx, state, err := parseNumberState(r)

		if err == 0 {
			control.FlagHTTPWaitResp = true
			control.SendObjectPwr(control.OBJECT_MODEM, idx, state)
			status, ret := waitForResponce(1)
			if ret {
				res.Number = idx + 1
				res.State = state
				resp.Results = &res
				control.PowerSt.Modem[idx] = state
			}
			resp.Status = status
		} else {
			resp.Status = "INVALID_REQUEST"
		}
	} else {
		resp.Status = "CHANEL_BUSY"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetDownPwrPC(w http.ResponseWriter, r *http.Request) {
	var res RespStateResults
	var resp RespState

	if control.USBChanWaitNotBusy(1000) {
		control.FlagHTTPWaitResp = true
		control.SendObjectPwr(control.OBJECT_PC, 0, false)
		status, ret := waitForResponce(1)
		if ret {
			res.Number = 0
			res.State = false
			resp.Results = &res
			control.PowerSt.Pc = false
		}
		resp.Status = status
	} else {
		resp.Status = "CHANEL_BUSY"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetPwrRelayByID(w http.ResponseWriter, r *http.Request) {
	var res RespStateResults
	var resp RespState

	if control.USBChanWaitNotBusy(1000) {
		idx, state, err := parseNumberState(r)

		if err == 0 {
			control.FlagHTTPWaitResp = true
			control.SendObjectPwr(control.OBJECT_RELAY, idx, state)
			status, ret := waitForResponce(1)
			if ret {
				res.Number = idx + 1
				res.State = state
				resp.Results = &res
				control.PowerSt.Relay[idx] = state
			}
			resp.Status = status
		} else {
			resp.Status = "INVALID_REQUEST"
		}
	} else {
		resp.Status = "CHANEL_BUSY"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetPwrWiFi(w http.ResponseWriter, r *http.Request) {
	var res RespStateResults
	var resp RespState

	if control.USBChanWaitNotBusy(1000) {
		_, state, err := parseNumberState(r)

		if err == 0 {
			control.FlagHTTPWaitResp = true
			control.SendObjectPwr(control.OBJECT_WIFI, 0, state)
			status, ret := waitForResponce(1)
			if ret {
				res.Number = 0
				res.State = state
				resp.Results = &res
				control.PowerSt.Wifi = state
			}
			resp.Status = status
		} else {
			resp.Status = "INVALID_REQUEST"
		}
	} else {
		resp.Status = "CHANEL_BUSY"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetWaitmode(w http.ResponseWriter, r *http.Request) {
	var res RespStateResults
	var resp RespState
	var ret bool
	var status string

	if control.USBChanWaitNotBusy(1000) {
		_, state, err := parseNumberState(r)

		if err == 0 {
			control.FlagHTTPWaitResp = true
			if state {
				control.SendCommand(control.CMD_PC_WAITMODE, true)
				status, ret = waitForResponce(1)
			} else {
				control.SendCommand(control.CMD_PC_SHUTDOWN, true)
				status, ret = waitForResponce(1)
			}

			if ret {
				res.Number = 0
				res.State = true
				resp.Results = &res
				control.PowerSt.Waitmode = false
			}
			resp.Status = status
		} else {
			resp.Status = "INVALID_REQUEST"
		}
	} else {
		resp.Status = "CHANEL_BUSY"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}
