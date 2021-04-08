package control

const (
	CMD_NONE           = 0
	CMD_LOCK           = 1  // [cmd, len, type]
	CMD_UNLOCK         = 2  // [cmd, len, type]
	CMD_FLYMODE        = 3  // [cmd, len, idx, state]
	CMD_AT_CMD         = 4  // [cmd, len, ]
	CMD_POWER          = 5  // [cmd, len, obj, idx, state]
	CMD_CHANGE_SIM     = 6  // [cmd, len, idx, num]
	CMD_LCD_PRINT      = 7  // [cmd, len, type, state]
	CMD_LCD_BLINK      = 8  // [cmd, len, idx, state]
	CMD_SET_IMEI       = 9  // [cmd, len, idx, data]
	CMD_SET_CONFIG     = 11 // [cmd, len, data]
	CMD_GET_CONFIG     = 12 // [cmd, len, data]
	CMD_MODEM_READY    = 13
	CMD_CFG_ERROR      = 14 // [cmd, len, 0]
	CMD_CTRL_ERROR     = 15 // [cmd, len, 0]
	CMD_PC_WAITMODE    = 16 // [cmd, len, 1]
	CMD_PC_SHUTDOWN    = 17 // [cmd, len, 1]
	CMD_PC_READY       = 18 // [cmd, len, 1]
	CMD_NEW_PHONES     = 19 // [cmd, len, data]
	CMD_SEND_SMS       = 20 // [cmd, len, idx, type, phone, msg]
	CMD_REQ_MODEM_INFO = 21 // [cmd, len, idx]
	CMD_REQ_CONN_INFO  = 22 // [cmd, len, idx]
	CMD_REQ_PHONES     = 23 // [cmd, len, 0]
	CMD_REQ_REASON     = 24 // [cmd, len, 0]
	CMD_OUT_SHUTDOWN   = 25 // [cmd, len, 1]
	CMD_OUT_SAVE_STATE = 26 // [cmd, len, data]
	CMD_OUT_SIM_CHANGE = 27 // [cmd, len, data]
	CMD_OUT_SMS        = 28 // [cmd, len, idx, type, phone, msg]
	CMD_OUT_AT_CMD     = 29 // [cmd, len, ]

	OBJECT_PC        = 1
	OBJECT_MODEM     = 2
	OBJECT_WIFI      = 3
	OBJECT_RELAY     = 4
	OBJECT_SMS_MODEM = 5

	IMEI_SIZE   = 15
	PHONE_SIZE  = 16
	ICCID_SIZE  = 18
	OPERID_SIZE = 5

	CONFIG_LEN  = 14
	PHONES_ROWS = 12
	PHONES_COLL = 3
)

type PowerStatus struct {
	PowerStat bool    // Static or battery power
	BatLevel  uint8   // Battery level
	Pc        bool    // PC power control
	Wifi      bool    // Wifi power control
	Relay     [2]bool // Relay power control
	Modem     [2]bool // Modem  power control
	Waitmode  bool    // Waitmode
}

type ModemStatus struct {
	Flymode bool   // Flightmode state
	SimNum  uint8  // Number of current sim-card in bank
	Iccid   string // ICCID of current sim-card
	Imei    string // IMEI of modem
	Phone   string // Current phone number
}

type ModemConnStatus struct {
	Status uint8
	Csq    uint8  // Signal level
	OperID string // Current operator //!Is operId the same as operator?
	CellID uint32 // ID of current base station
	Tac    uint16
}

type SystemStatus struct {
	SmsLock     bool
	ButtonsLock bool
	ReasonBuf   []byte
}

type ModemPhones struct {
	PhonesIn  [4]string
	PhonesOut [4]string
}

type ModemSimParams struct {
	Iccid  string
	Imei   string
	OperID string
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
