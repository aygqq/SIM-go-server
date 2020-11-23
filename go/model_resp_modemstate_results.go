/*
 * Power control block
 *
 * This API was created to monitor states of Power Control Block and send some commands to it.
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RespModemstateResults struct {
	// Number
	Number uint8 `json:"number"`
	// Flightmode state
	Flymode bool `json:"flymode"`
	// Number of current sim-card in bank
	SimNum uint8 `json:"sim_num"`
	// ICCID of current sim-card
	Imsi string `json:"sim_id"`
	// IMEI of modem
	Imei string `json:"imei"`
	// Current phone number
	Phone string `json:"phone"`
}
