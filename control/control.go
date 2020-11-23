package control

import (
	"container/list"
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

var PowerSt PowerStatus
var ModemSt [2]ModemStatus
var ConnSt [2]ModemConnStatus
var SystemSt SystemStatus
var ModemPh ModemPhones

var SmsList *list.List

var modemPhReq ModemPhones
var modemStReq ModemStatus

var HttpReqChan chan uint8 = make(chan uint8)
var ControlReqChan chan uint8 = make(chan uint8)
var FlagWaitResp bool = false

func waitForResponce() error {
	var err error

	select {
	case read := <-ControlReqChan:
		//! COM now in echo mode, so that "read" value doesn't matter
		// if read == 0 {
		// 	err...
		// }
		fmt.Printf("Chanel recv %d\n", read)
	case <-time.After(time.Second):
		fmt.Println("No response received")
		err = errors.New("No response received")
	}
	return err
}

func Init() error {
	ph, err := readPhonesFile("../phones.csv")
	if err != nil {
		fmt.Printf("Failed to read file: %q", err)
		SendCommand(CMD_CFG_ERROR, true)
		waitForResponce()
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
		if reflect.DeepEqual(modemPhReq, ph.Phones) == false {
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
		ProcLactConfigStart()
	} else if strings.HasPrefix(SystemSt.ReasonBuf, "Modem") {
		var cfg ModemPowerConfig
		data := []byte(SystemSt.ReasonBuf)
		cfg.m1Pwr = data[6]
		cfg.m1Sim = data[7]
		cfg.m2Pwr = data[8]
		cfg.m2Sim = data[9]
		ProcModemStart(&cfg)
	}

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
	if PowerSt.Modem[idx] == true {
		SendObjectPwr(OBJECT_MODEM, idx, false)
		if err = waitForResponce(); err != nil {
			return err
		}
		time.Sleep(10 * time.Second)
		SendObjectPwr(OBJECT_MODEM, idx, true)
		if err = waitForResponce(); err != nil {
			return err
		}
	}

	SendFlightmode(idx, true)
	if err = waitForResponce(); err != nil {
		return err
	}
	//? The power is already turned on
	SendObjectPwr(OBJECT_MODEM, idx, true)
	if err = waitForResponce(); err != nil {
		return err
	}

	SendDoubleByte(CMD_CHANGE_SIM, idx, sim)
	if err = waitForResponce(); err != nil {
		return err
	}

	SendDoubleByte(CMD_LCD_BLINK, idx, 0)
	if err = waitForResponce(); err != nil {
		return err
	}

	//? How to wait for modem loaded?
	time.Sleep(10 * time.Second)

	SendShort(CMD_REQ_MODEM_INFO, idx)
	if err = waitForResponce(); err != nil {
		return err
	}
	//? Imsi should be changed on PCB by reading it from SIM?
	if modemStReq.Imsi != phFile.Bank[idx][sim].Imsi {
		fmt.Printf("Imsi is wrong")
		err = errors.New("Imsi is wrong")
		SendCommand(CMD_CFG_ERROR, true)
		waitForResponce()
		return err
	}
	if modemStReq.Imei != phFile.Bank[idx][sim].Imei {
		SendSetImei(idx, phFile.Bank[idx][sim].Imei)
		if err = waitForResponce(); err != nil {
			return err
		}

		SendShort(CMD_REQ_MODEM_INFO, idx)
		if err = waitForResponce(); err != nil {
			return err
		}

		if modemStReq.Imei != phFile.Bank[idx][sim].Imei {
			fmt.Printf("Can not set IMEI")
			err = errors.New("Can not set IMEI")
			SendCommand(CMD_CFG_ERROR, true)
			waitForResponce()
			return err
		}
	}

	SendFlightmode(idx, false)
	if err = waitForResponce(); err != nil {
		return err
	}

	return nil
}

func ProcButtonStart() {

}

func ProcSetConfigStart() {
	//!Glubokiy shit
}

func ProcLactConfigStart() error {
	cfg, err := readConfigFile("../config.txt")
	if err != nil {
		fmt.Printf("Failed to read file: %q", err)
		SendCommand(CMD_CTRL_ERROR, true)
		waitForResponce()
		SendCommand(CMD_PC_SHUTDOWN, true)
		waitForResponce()
		return err
	}

	SendConfig(cfg)

	return nil
}

func ProcModemStart(cfg *ModemPowerConfig) {
	var err error

	if cfg.m2Pwr == 1 {
		err = modemTurnOn(1, cfg.m2Sim)
		if err != nil {
			fmt.Printf("Failed to turn on modem 2: %q", err)
		}
	}
	if cfg.m1Pwr == 1 {
		err = modemTurnOn(0, cfg.m1Sim)
		if err != nil {
			fmt.Printf("Failed to turn on modem 1: %q", err)
		}
	}

	SendShort(CMD_UNLOCK, 1)
	if err = waitForResponce(); err != nil {
		return
	}
}
