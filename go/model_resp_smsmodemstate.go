/*
 * API для взаимодействия с КУЭП
 *
 * Данное API сделано с целью мониторинга состояний контроллера управления электропитанием (КУЭП), а также для отправки КУЭП различных команд.
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RespSmsmodemstate struct {

	Results *RespSmsmodemstateResults `json:"results,omitempty"`
	// Three possible statuses:   * `OK`: No errors occurred.  * `INVALID_REQUEST`: Some parameters are missing or invalid.  * `EXECUTE_ERROR`: No or wrong responce from Power Control Block.  * `UNKNOWN_ERROR`: The request could not be processed due to a server error. The request may succeed if you try again.
	Status string `json:"status,omitempty"`
}
