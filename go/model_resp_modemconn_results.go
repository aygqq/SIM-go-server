/*
 * Power control block
 *
 * This API was created to monitor states of Power Control Block and send some commands to it.
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RespModemconnResults struct {
	// Number
	Number uint8 `json:"number"`
	// Current operator
	Operator string `json:"operator,omitempty"`
	// ID of current base station
	BaseID string `json:"base_id,omitempty"`
	// Signal level
	Signal string `json:"signal,omitempty"`
}
