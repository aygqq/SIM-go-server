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
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var HttpReqChan chan uint8 = make(chan uint8)
var FlagWaitResp bool = false

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// func MakeChanel() *chan uint8 {
// 	RequestChan = make(chan uint8)
// 	return &RequestChan
// }

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"GetPwrBat",
		strings.ToUpper("Get"),
		"/power/battery",
		GetPwrBat,
	},

	Route{
		"SetButtonsLock",
		strings.ToUpper("Put"),
		"/buttons/lock",
		SetButtonsLock,
	},

	Route{
		"GetFileConfig",
		strings.ToUpper("Get"),
		"/files/config",
		GetFileConfig,
	},

	Route{
		"SetFileConfig",
		strings.ToUpper("Put"),
		"/files/config",
		SetFileConfig,
	},

	Route{
		"GetFilePhones",
		strings.ToUpper("Get"),
		"/files/phones",
		GetFilePhones,
	},

	Route{
		"SetFileNPhones",
		strings.ToUpper("Put"),
		"/files/phones",
		SetFileNPhones,
	},

	Route{
		"GetModemConnByID",
		strings.ToUpper("Get"),
		"/modem/conn/[modem_id]",
		GetModemConnByID,
	},

	Route{
		"GetModemFlyByID",
		strings.ToUpper("Get"),
		"/modem/state/flymode/[modem_id]",
		GetModemFlyByID,
	},

	Route{
		"SetModemFlyByID",
		strings.ToUpper("Put"),
		"/modem/state/flymode/[modem_id]",
		SetModemFlyByID,
	},

	Route{
		"GetModemImeiByID",
		strings.ToUpper("Get"),
		"/modem/state/imei/[modem_id]",
		GetModemImeiByID,
	},

	Route{
		"SetModemImeiByID",
		strings.ToUpper("Put"),
		"/modem/state/imei/[modem_id]",
		SetModemImeiByID,
	},

	Route{
		"GetModemSimByID",
		strings.ToUpper("Get"),
		"/modem/state/sim/[modem_id]",
		GetModemSimByID,
	},

	Route{
		"SetModemSimByID",
		strings.ToUpper("Put"),
		"/modem/state/sim/[modem_id]",
		SetModemSimByID,
	},

	Route{
		"GetModemStByID",
		strings.ToUpper("Get"),
		"/modem/state/[modem_id]",
		GetModemStByID,
	},

	Route{
		"GetPwrCfg",
		strings.ToUpper("Get"),
		"/power",
		GetPwrCfg,
	},

	Route{
		"SetPwrCfg",
		strings.ToUpper("Put"),
		"/power",
		SetPwrCfg,
	},

	Route{
		"GetPwrModemByID",
		strings.ToUpper("Get"),
		"/power/modem/[modem_id]",
		GetPwrModemByID,
	},

	Route{
		"SetPwrModemByID",
		strings.ToUpper("Put"),
		"/power/modem/[modem_id]",
		SetPwrModemByID,
	},

	Route{
		"GetPwrPC",
		strings.ToUpper("Get"),
		"/power/pc",
		GetPwrPC,
	},

	Route{
		"SetPwrPC",
		strings.ToUpper("Put"),
		"/power/pc",
		SetPwrPC,
	},

	Route{
		"GetPwrRelayByID",
		strings.ToUpper("Get"),
		"/power/relay/[rel_id]",
		GetPwrRelayByID,
	},

	Route{
		"SetPwrRelayByID",
		strings.ToUpper("Put"),
		"/power/relay/[rel_id]",
		SetPwrRelayByID,
	},

	Route{
		"GetPwrWiFi",
		strings.ToUpper("Get"),
		"/power/wifi",
		GetPwrWiFi,
	},

	Route{
		"SetPwrWiFi",
		strings.ToUpper("Put"),
		"/power/wifi",
		SetPwrWiFi,
	},

	Route{
		"SetSmsLock",
		strings.ToUpper("Put"),
		"/sms/lock",
		SetSmsLock,
	},
}
