package control

const (
	CMD_NONE           = 0
	CMD_LOCK           = 1  // [cmd, len, type]
	CMD_UNLOCK         = 2  // [cmd, len, 0]
	CMD_FLYMODE        = 3  // [cmd, len, idx, state]
	CMD_AT_CMD         = 4  // [cmd, len, ]
	CMD_POWER          = 5  // [cmd, len, obj, idx, state]
	CMD_CHANGE_SIM     = 6  // [cmd, len, idx, num]
	CMD_LCD_PRINT      = 7  // [cmd, len, type, state]
	CMD_LCD_BLINK      = 8  // [cmd, len, idx, state]
	CMD_SET_IMEI       = 9  // [cmd, len, idx, data]
	CMD_SET_CONFIG     = 10 // [cmd, len, data]
	CMD_CFG_ERROR      = 11 // [cmd, len, 0]
	CMD_CTRL_ERROR     = 12 // [cmd, len, 0]
	CMD_PC_WAITMODE    = 13 // [cmd, len, 1]
	CMD_PC_SHUTDOWN    = 14 // [cmd, len, 1]
	CMD_PC_READY       = 15 // [cmd, len, 1]
	CMD_NEW_PHONES     = 16 // [cmd, len, data]
	CMD_SEND_SMS       = 17 // [cmd, len, idx, type, phone, msg]
	CMD_REQ_MODEM_INFO = 18 // [cmd, len, idx]
	CMD_REQ_PHONES     = 19 // [cmd, len, 0]
	CMD_REQ_REASON     = 20 // [cmd, len, 0]
	CMD_OUT_SHUTDOWN   = 21 // [cmd, len, 1]
	CMD_OUT_SAVE_STATE = 22 // [cmd, len, data]
	CMD_OUT_SIM_CHANGE = 23 // [cmd, len, data]
	CMD_OUT_SMS        = 24 // [cmd, len, idx, type, phone, msg]
	CMD_OUT_AT_CMD     = 25 // [cmd, len, ]

	OBJECT_PC        = 1
	OBJECT_MODEM     = 2
	OBJECT_WIFI      = 3
	OBJECT_RELAY     = 4
	OBJECT_SMS_MODEM = 5

	IMEI_SIZE   = 15
	PHONE_SIZE  = 26
	IMSI_SIZE   = 15
	OPERID_SIZE = 5

	CONFIG_LEN  = 14
	PHONES_ROWS = 12
	PHONES_COLL = 3
)

type PowerStatus struct {
	// Static or battery power
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
	Imsi string
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

type SystemStatus struct {
	SmsLock     bool
	ButtonsLock bool
	ReasonBuf   string
}

type ModemPhones struct {
	PhonesIn  [4]string
	PhonesOut [4]string
}

type ModemSimParams struct {
	Imsi   string
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
	Bank   [2][4]ModemSimParams
	Phones ModemPhones
}

type ModemPowerConfig struct {
	m1Pwr uint8
	m1Sim uint8
	m2Pwr uint8
	m2Sim uint8
}

type SmsMessage struct {
	ModemNum uint8
	MsgType  uint8
	Phone    string
	Message  string
}
