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
	"io/ioutil"
	"log"
	"net/http"

	"../control"
)

func GetFileConfig(w http.ResponseWriter, r *http.Request) {
	var resp RespFilecfg
	resp.Results = control.GetConfigFileString()
	resp.Status = "OK"

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func GetFilePhones(w http.ResponseWriter, r *http.Request) {
	var resp RespFilephones

	control.GetPhonesFile(&resp.Results)
	resp.Status = "OK"

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetFileConfig(w http.ResponseWriter, r *http.Request) {
	var resp RespFilecfg
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}
	cfg := string(body)
	//log.Printf("Request body is %s\n", cfg)

	control.SetConfigFile(cfg)

	resp.Results = control.GetConfigFileString()
	resp.Status = "OK"

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	str, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(str))
}

func SetFileNPhones(w http.ResponseWriter, r *http.Request) {
	var resp RespFilephones

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}
	//str := string(body)
	//log.Printf("Request body is %s\n", str)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	err = json.Unmarshal(body, &resp.Results)
	control.SetPhonesFile(&resp.Results)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
