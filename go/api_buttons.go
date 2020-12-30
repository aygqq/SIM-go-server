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

	"../control"
)

func SetButtonsLock(w http.ResponseWriter, r *http.Request) {
	var res RespStateResults
	var resp RespState

	_, state, err := parseNumberState(r)
	if err == 0 {
		if state == true {
			control.SendShort(control.CMD_LOCK, 2)
		} else {
			control.SendShort(control.CMD_UNLOCK, 2)
		}
		status, ret := waitForResponce(1)
		if ret == true {
			res.Number = 0
			res.State = state
			resp.Results = &res
			control.SystemSt.ButtonsLock = state
		}
		resp.Status = status
	} else {
		resp.Status = "INVALID_REQUEST"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}
