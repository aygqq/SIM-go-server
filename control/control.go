package control

import (
	"container/list"
	"errors"
	"log"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

// PowerSt - States of power control block
var PowerSt PowerStatus

// ModemSt - States of modems
var ModemSt [2]ModemStatus

// ConnSt - Modem connection states
var ConnSt [2]ModemConnStatus

// SystemSt - Some system states
var SystemSt SystemStatus

// ModemPh -
var ModemPh ModemPhones

// SmsList - List of recieved sms messages
var SmsList *list.List

var modemPhReq ModemPhones
var modemStReq ModemStatus

// HTTPReqChan - Chanel to proceed reply to API
var HTTPReqChan chan uint8 = make(chan uint8)

// ControlReqChan - Chanel to proceed reply to control
var ControlReqChan chan uint8 = make(chan uint8)

// FlagHTTPWaitResp - What chanel is in use
var FlagHTTPWaitResp bool = false

// FlagControlWaitResp - What chanel is in use
var FlagControlWaitResp bool = false

func waitForResponce() error {
	var err error

	FlagControlWaitResp = true

	select {
	case read := <-ControlReqChan:
		//! COM now in echo mode, so that "read" value doesn't matter
		if read == 0 {
			err = errors.New("Wrong response received")
		}
		// log.Printf("Chanel recv %d\n", read)
	case <-time.After(2 * time.Second):
		log.Println("No response received")
		err = errors.New("No response received")
	}

	FlagControlWaitResp = false

	return err
}

// ProcStart function
func ProcStart() error {
	ph, err := readPhonesFile("phones.csv")
	if err != nil {
		log.Printf("Failed to read file: %q\n", err)
		SendCommand(CMD_CFG_ERROR, true)
		waitForResponce()
		return err
	}

	log.Println("\tFlightmode on")
	SendFlightmode(0, true)
	if err = waitForResponce(); err != nil {
		return err
	}

	SendCommand(CMD_REQ_PHONES, true)
	if err = waitForResponce(); err != nil {
		return err
	}

	if reflect.DeepEqual(modemPhReq, ph.Phones) == false {
		SendNewPhones(ph.Phones)
		if err = waitForResponce(); err != nil {
			return err
		}
		SendCommand(CMD_REQ_PHONES, true)
		if err = waitForResponce(); err != nil {
			return err
		}
		if modemPhReq == ph.Phones {
			log.Println("Phones are equal")
		} else {
			log.Println("Phones recv\n", modemPhReq)
			log.Println("Phones file\n", ph.Phones)
			err = errors.New("Phones file double check failed")
			SendCommand(CMD_CFG_ERROR, true)
			waitForResponce()
			return err
		}
	}

	SendCommand(CMD_PC_READY, true)
	if err = waitForResponce(); err != nil {
		return err
	}

	SendCommand(CMD_REQ_REASON, true)
	if err = waitForResponce(); err != nil {
		return err
	}
	log.Printf("Reason buf is %s\n", SystemSt.ReasonBuf)
	if strings.HasPrefix(SystemSt.ReasonBuf, "Button") {
		ProcButtonStart()
	} else if strings.HasPrefix(SystemSt.ReasonBuf, "Sms") {
		ProcSetConfigStart()
	} else if strings.HasPrefix(SystemSt.ReasonBuf, "Last") {
		ProcLastConfigStart()
	} else if strings.HasPrefix(SystemSt.ReasonBuf, "Modem") {
		var cfg ModemPowerConfig
		data := []byte(SystemSt.ReasonBuf)
		cfg.m1Pwr = data[6]
		cfg.m1Sim = data[7]
		cfg.m2Pwr = data[8]
		cfg.m2Sim = data[9]
		ProcModemStart(&cfg)
	}

	phFile = ph

	return nil
}

func procShutdown() {
	err := exec.Command("/bin/sh", "shutdown.sh").Run()
	if err != nil {
		log.Println(err)
	}
}

func modemTurnOn(idx uint8, sim uint8) error {
	var err error
	log.Printf("Modem %d turn on\n", idx+1)
	if PowerSt.Modem[idx] == true {
		log.Println("\tPower off")
		SendObjectPwr(OBJECT_MODEM, idx, false)
		if err = waitForResponce(); err != nil {
			return err
		}
		time.Sleep(10 * time.Second)
	}

	log.Println("\tPower on")
	SendObjectPwr(OBJECT_MODEM, idx, true)
	if err = waitForResponce(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	log.Println("\tFlightmode on")
	SendFlightmode(idx, true)
	if err = waitForResponce(); err != nil {
		return err
	}

	log.Println("\tChange sim")
	SendDoubleByte(CMD_CHANGE_SIM, idx, sim)
	if err = waitForResponce(); err != nil {
		return err
	}

	log.Println("\tLCD blink")
	SendDoubleByte(CMD_LCD_BLINK, idx, 0)
	if err = waitForResponce(); err != nil {
		return err
	}

	//? How to wait for modem loaded?
	time.Sleep(35 * time.Second)

	log.Println("\tReq modem info")
	SendShort(CMD_REQ_MODEM_INFO, idx)
	if err = waitForResponce(); err != nil {
		return err
	}
	//? Iccid should be changed on PCB by reading it from SIM?
	if modemStReq.Iccid != phFile.Bank[idx][sim-1].Iccid {
		log.Printf("\tIccid is wrong %s %s\n", modemStReq.Iccid, phFile.Bank[idx][sim-1].Iccid)
		err = errors.New("Iccid is wrong")
		SendCommand(CMD_CFG_ERROR, true)
		waitForResponce()
		return err
	}
	log.Printf("\tIccid is %s\n", modemStReq.Iccid)
	if modemStReq.Imei != phFile.Bank[idx][sim-1].Imei {
		time.Sleep(5 * time.Second)
		SendSetImei(idx, phFile.Bank[idx][sim-1].Imei)
		if err = waitForResponce(); err != nil {
			return err
		}

		SendShort(CMD_REQ_MODEM_INFO, idx)
		if err = waitForResponce(); err != nil {
			return err
		}

		if modemStReq.Imei != phFile.Bank[idx][sim-1].Imei {
			log.Printf("\tCan not set IMEI %s %s\n", modemStReq.Imei, phFile.Bank[idx][sim-1].Imei)
			err = errors.New("Can not set IMEI")
			SendCommand(CMD_CFG_ERROR, true)
			waitForResponce()
			return err
		}
	}
	log.Printf("\tIMEI is %s\n", modemStReq.Imei)

	log.Println("\tFlightmode off")
	SendFlightmode(idx, false)
	if err = waitForResponce(); err != nil {
		return err
	}
	log.Printf("Modem %d turn on (END)\n", idx+1)

	return nil
}

// ProcButtonStart -
func ProcButtonStart() {

}

// ProcSetConfigStart -
func ProcSetConfigStart() {

}

// ProcLastConfigStart - Work on last config
func ProcLastConfigStart() error {
	cfg, err := readConfigFile("config.txt")
	if err != nil {
		log.Printf("Failed to read file: %q\n", err)
		SendCommand(CMD_CTRL_ERROR, true)
		waitForResponce()
		SendCommand(CMD_PC_SHUTDOWN, true)
		waitForResponce()
		return err
	}

	SendConfig(cfg)

	return nil
}

// ProcModemStart - This procedure sterts the modem
func ProcModemStart(cfg *ModemPowerConfig) {
	var err error

	if cfg.m2Pwr == 1 {
		err = modemTurnOn(1, cfg.m2Sim)
		if err != nil {
			log.Printf("Failed to turn on modem 2: %q\n", err)

			ModemSt[1].Iccid = ""
			ModemSt[1].Imei = ""
			ModemSt[1].Flymode = false
			ModemSt[1].SimNum = 0
			ModemSt[1].Phone = ""

			SendCommand(CMD_CFG_ERROR, true)
			waitForResponce()
		}
	}
	if cfg.m1Pwr == 1 {
		err = modemTurnOn(0, cfg.m1Sim)
		if err != nil {
			log.Printf("Failed to turn on modem 1: %q\n", err)

			ModemSt[0].Iccid = ""
			ModemSt[0].Imei = ""
			ModemSt[0].Flymode = false
			ModemSt[0].SimNum = 0
			ModemSt[0].Phone = ""

			SendCommand(CMD_CFG_ERROR, true)
			waitForResponce()
		}
	}

	SendShort(CMD_UNLOCK, 1)
	if err = waitForResponce(); err != nil {
		log.Printf("Cmd unlock: %q\n", err)
		return
	}
}
