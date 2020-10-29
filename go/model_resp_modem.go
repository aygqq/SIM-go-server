/*
 * Power control block
 *
 * This API was created to monitor states of Power Control Block and send some commands to it.
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RespModem struct {
	// State
	Results bool `json:"results,omitempty"`
	// Three possible statuses:   * `OK`: No errors occurred.  * `UNKNOWN_ERROR`: The request could not be processed due to a server error. The request may succeed if you try again.
	Status string `json:"status,omitempty"`
}
