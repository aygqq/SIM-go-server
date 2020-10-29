package control

const (
	CMD_NONE           = 0
	CMD_LOCK           = 1
	CMD_UNLOCK         = 2
	CMD_FLYMODE        = 3
	CMD_AT_CMD         = 4
	CMD_POWER          = 5
	CMD_CHANGE_SIM     = 6
	CMD_LCD_PRINT      = 7
	CMD_SET_IMEI       = 8
	CMD_SET_CONFIG     = 9
	CMD_CFG_ERROR      = 10
	CMD_CTRL_ERROR     = 11
	CMD_PC_WAITMODE    = 12
	CMD_PC_SHUTDOWN    = 13
	CMD_PC_READY       = 14
	CMD_NEW_PHONES     = 15
	CMD_REQ_MODEM_INFO = 16
	CMD_REQ_PHONES     = 17
	CMD_REQ_REASON     = 18
	CMD_OUT_SHUTDOWN   = 19
	CMD_OUT_SAVE_STATE = 20
	CMD_OUT_SIM_CHANGE = 21
	CMD_OUT_AT_CMD     = 22

	IMEI_SIZE  = 16
	PHONE_SIZE = 26
	SIMID_SIZE = 20
)

type PowerStatus struct {
	// Static power or battery
	PowerStat bool
	// Battery level
	BatLevel int32
	// PC power control
	Pc bool
	// Wifi power control
	Wifi bool
	// Relay power control
	Relay [2]bool
	// Modem  power control
	Modem [2]bool
}

type ModemStatus struct {
	// Flightmode state
	Flymode bool
	// Number of current sim-card in bank
	SimNum uint8
	// ICCID of current sim-card
	SimId string
	// IMEI of modem
	Imei string
	// Current phone number
	Phone string
}

type ModemConnStatus struct {
	// Current operator
	Operator string
	// ID of current base station
	BaseId string
	// Signal level
	Signal string
}

type RequiredElement struct {
	State  *bool
	Number *int32
	String *string
}

type SystemStatus struct {
	SmsLock     bool
	ButtonsLock bool
	ReasonBuf   []byte
	ReqElem     RequiredElement
}

type ModemPhones struct {
	phonesIn  [4]string
	phonesOut [4]string
}
