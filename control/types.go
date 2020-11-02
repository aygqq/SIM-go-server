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
	powerStat bool
	// Battery level
	batLevel uint8
	// PC power control
	pc bool
	// Wifi power control
	wifi bool
	// Relay power control
	relay [2]bool
	// Modem  power control
	modem [2]bool
}

type ModemStatus struct {
	// Flightmode state
	flymode bool
	// Number of current sim-card in bank
	simNum uint8
	// ICCID of current sim-card
	simId string
	// IMEI of modem
	imei string
	// Current phone number
	phone string
}

type ModemConnStatus struct {
	// Current operator
	operator string //!Is operId the same as operator?
	// ID of current base station
	baseId string
	// Signal level
	signal string
}

type RequiredElement struct {
	State  *bool
	Number *int32
	String *string
}

type SystemStatus struct {
	smsLock     bool
	buttonsLock bool
	reasonBuf   []byte
	reqElem     RequiredElement
}

type ModemPhones struct {
	phonesIn  [4]string
	phonesOut [4]string
}

type ModemSimParams struct {
	simId  string
	imei   string
	operId string //!Is operId the same as operator?
}

type FileConfig struct {
	power      PowerStatus
	simNum     [2]uint8
	configErr  bool
	stateErr   bool
	connectErr bool
}

type FilePhones struct {
	bank1  [4]ModemSimParams
	bank2  [4]ModemSimParams
	phones ModemPhones
}
