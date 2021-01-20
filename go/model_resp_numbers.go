/*
 * Power control block
 *
 * This API was created to monitor states of Power Control Block and send some commands to it.
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RespNumbers struct {
	// Разрешенные телефонные номера
	Results [4][2]string `json:"results,omitempty"`
	// Three possible statuses:   * `OK`: No errors occurred.  * `INVALID_REQUEST`: Some parameters are missing or invalid.  * `EXECUTE_ERROR`: No or wrong responce from Power Control Block.  * `UNKNOWN_ERROR`: The request could not be processed due to a server error. The request may succeed if you try again.
	Status string `json:"status,omitempty"`
}