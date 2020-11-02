package control

var PowerSt PowerStatus
var ModemSt [2]ModemStatus
var ConnSt [2]ModemConnStatus
var SystemSt SystemStatus
var ModemPh ModemPhones

var HttpReqChan chan uint8 = make(chan uint8)
var ControlReqChan chan uint8 = make(chan uint8)
var FlagWaitResp bool = false

func GetPowerConfig() PowerStatus {
	return PowerSt
}

func SetPowerConfig(cfg PowerStatus) {

}
