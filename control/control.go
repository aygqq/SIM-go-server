package control

import (
	"container/list"
	"errors"
	"fmt"
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
		log.Printf("Chanel recv %d\n", read)
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
		fmt.Printf("Failed to read file: %q\n", err)
		SendCommand(CMD_CFG_ERROR, true)
		waitForResponce()
		return err
	}

	fmt.Println("\tFlightmode on")
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
			fmt.Println("GOOOOOOOD")
		} else {
			fmt.Println("Phones recv\n", modemPhReq)
			fmt.Println("Phones file\n", ph.Phones)
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
	err := exec.Command("/bin/sh", "/app/shutdown.sh").Run()
	if err != nil {
		fmt.Println(err)
	}
}

func modemTurnOn(idx uint8, sim uint8) error {
	var err error
	fmt.Printf("Modem %d turn on\n", idx+1)
	if PowerSt.Modem[idx] == true {
		fmt.Println("\tPower off")
		SendObjectPwr(OBJECT_MODEM, idx, false)
		if err = waitForResponce(); err != nil {
			return err
		}
		time.Sleep(10 * time.Second)
	}

	fmt.Println("\tPower on")
	SendObjectPwr(OBJECT_MODEM, idx, true)
	if err = waitForResponce(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	fmt.Println("\tFlightmode on")
	SendFlightmode(idx, true)
	if err = waitForResponce(); err != nil {
		return err
	}

	fmt.Println("\tChange sim")
	SendDoubleByte(CMD_CHANGE_SIM, idx, sim)
	if err = waitForResponce(); err != nil {
		return err
	}

	fmt.Println("\tLCD blink")
	SendDoubleByte(CMD_LCD_BLINK, idx, 0)
	if err = waitForResponce(); err != nil {
		return err
	}

	//? How to wait for modem loaded?
	time.Sleep(35 * time.Second)

	fmt.Println("\tReq modem info")
	SendShort(CMD_REQ_MODEM_INFO, idx)
	if err = waitForResponce(); err != nil {
		return err
	}
	//? Iccid should be changed on PCB by reading it from SIM?
	if modemStReq.Iccid != phFile.Bank[idx][sim-1].Iccid {
		fmt.Printf("\tIccid is wrong %s %s\n", modemStReq.Iccid, phFile.Bank[idx][sim-1].Iccid)
		err = errors.New("Iccid is wrong")
		SendCommand(CMD_CFG_ERROR, true)
		waitForResponce()
		return err
	}
	if modemStReq.Imei != phFile.Bank[idx][sim-1].Imei {
		SendSetImei(idx, phFile.Bank[idx][sim-1].Imei)
		if err = waitForResponce(); err != nil {
			return err
		}

		SendShort(CMD_REQ_MODEM_INFO, idx)
		if err = waitForResponce(); err != nil {
			return err
		}

		if modemStReq.Imei != phFile.Bank[idx][sim-1].Imei {
			fmt.Printf("\tCan not set IMEI %s %s\n", modemStReq.Imei, phFile.Bank[idx][sim-1].Imei)
			err = errors.New("Can not set IMEI")
			SendCommand(CMD_CFG_ERROR, true)
			waitForResponce()
			return err
		}
	}

	fmt.Println("\tFlightmode off")
	SendFlightmode(idx, false)
	if err = waitForResponce(); err != nil {
		return err
	}
	fmt.Printf("Modem %d turn on (END)\n", idx+1)

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
	cfg, err := readConfigFile("../config.txt")
	if err != nil {
		fmt.Printf("Failed to read file: %q\n", err)
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
			fmt.Printf("Failed to turn on modem 2: %q\n", err)
		}
	}
	if cfg.m1Pwr == 1 {
		err = modemTurnOn(0, cfg.m1Sim)
		if err != nil {
			fmt.Printf("Failed to turn on modem 1: %q\n", err)
		}
	}

	SendShort(CMD_UNLOCK, 1)
	if err = waitForResponce(); err != nil {
		fmt.Printf("Cmd unlock: %q\n", err)
		return
	}
}
