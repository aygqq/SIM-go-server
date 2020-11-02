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
	CMD_LCD_BLINK      = 8
	CMD_SET_IMEI       = 9
	CMD_SET_CONFIG     = 10
	CMD_CFG_ERROR      = 11
	CMD_CTRL_ERROR     = 12
	CMD_PC_WAITMODE    = 13
	CMD_PC_SHUTDOWN    = 14
	CMD_PC_READY       = 15
	CMD_NEW_PHONES     = 16
	CMD_REQ_MODEM_INFO = 17
	CMD_REQ_PHONES     = 18
	CMD_REQ_REASON     = 19
	CMD_OUT_SHUTDOWN   = 20
	CMD_OUT_SAVE_STATE = 21
	CMD_OUT_SIM_CHANGE = 22
	CMD_OUT_AT_CMD     = 23

	OBJECT_PC        = 1
	OBJECT_MODEM     = 2
	OBJECT_WIFI      = 3
	OBJECT_RELAY     = 4
	OBJECT_SMS_MODEM = 5

	IMEI_SIZE  = 16
	PHONE_SIZE = 26
	SIMID_SIZE = 20
)

type PowerStatus struct {
	// Static power or battery
	PowerStat bool
	// Battery level
	BatLevel uint8
	// PC power control
	Pc bool
	// Wifi power control
	Wifi bool
	// Relay power control
	Relay [2]bool
	// Modem  power control
	Modem [2]bool
	// Waitmode
	Waitmode bool
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
	Operator string //!Is operId the same as operator?
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
	PhonesIn  [4]string
	PhonesOut [4]string
}

type ModemSimParams struct {
	SimId  string
	Imei   string
	OperId string //!Is operId the same as operator?
}

type FileConfig struct {
	Power      PowerStatus
	SimNum     [2]uint8
	ConfigErr  bool
	StateErr   bool
	ConnectErr bool
}

type FilePhones struct {
	Bank1  [4]ModemSimParams
	Bank2  [4]ModemSimParams
	Phones ModemPhones
}
