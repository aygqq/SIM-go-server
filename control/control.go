package control

import (
	"container/list"
	"errors"
	"flag"
	"log"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"gopkg.in/routeros.v2"
)

var (
	address  = flag.String("address", "192.168.88.1:8728", "Address")
	username = flag.String("username", "admin", "Username")
	password = flag.String("password", "", "Password")
)

// PowerSt - States of power control block
var PowerSt PowerStatus

// ModemSt - States of modems
var ModemSt [2]ModemStatus

// ConnSt - Modem connection states
var ConnSt [2]ModemConnStatus

// ModemSt - States of modems
var SmsModemSt [2]ModemStatus

// SmsConnSt - Modem connection states
var SmsConnSt [2]ModemConnStatus

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

var router [2]routerInfo

func waitForResponce() error {
	// FlagControlWaitResp = true

	var err error

	select {
	case read := <-ControlReqChan:
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

func ProcSetPhones(ph ModemPhones) error {
	var err error

	FlagControlWaitResp = true
	SendCommand(CMD_REQ_PHONES, true)
	if err = waitForResponce(); err != nil {
		return err
	}

	if reflect.DeepEqual(modemPhReq, ph) == false {
		FlagControlWaitResp = true
		SendNewPhones(ph)
		if err = waitForResponce(); err != nil {
			return err
		}
		FlagControlWaitResp = true
		SendCommand(CMD_REQ_PHONES, true)
		if err = waitForResponce(); err != nil {
			return err
		}
		if modemPhReq == ph {
			log.Println("Phones are equal")
		} else {
			log.Println("Phones recv\n", modemPhReq)
			log.Println("Phones file\n", ph)
			err = errors.New("Phones file double check failed")
			FlagControlWaitResp = true
			SendCommand(CMD_CFG_ERROR, true)
			waitForResponce()
			return err
		}
	}

	return nil
}

// ProcStart function
func ProcStart() error {
	err := readRouterFile("routers.csv")
	if err != nil {
		log.Printf("Failed to read file: %q\n", err)
		FlagControlWaitResp = true
		SendCommand(CMD_CFG_ERROR, true)
		waitForResponce()
		FlagControlWaitResp = true
		SendCommand(CMD_PC_READY, true)
		waitForResponce()
		return err
	}

	ph, err := readPhonesFile("phones.csv")
	if err != nil {
		log.Printf("Failed to read file: %q\n", err)
		FlagControlWaitResp = true
		SendCommand(CMD_CFG_ERROR, true)
		waitForResponce()
		FlagControlWaitResp = true
		SendCommand(CMD_PC_READY, true)
		waitForResponce()
		return err
	}

	err = checkPhonesFile(&ph)
	if err != nil {
		log.Printf("Failed to read file: %q\n", err)
		FlagControlWaitResp = true
		SendCommand(CMD_CFG_ERROR, true)
		waitForResponce()
		FlagControlWaitResp = true
		SendCommand(CMD_PC_READY, true)
		waitForResponce()
		return err
	}
	phFile = ph

	err = ProcSetPhones(ph.Phones)
	if err != nil {
		FlagControlWaitResp = true
		SendCommand(CMD_PC_READY, true)
		waitForResponce()
		return err
	}

	FlagControlWaitResp = true
	SendCommand(CMD_PC_READY, true)
	if err = waitForResponce(); err != nil {
		return err
	}

	FlagControlWaitResp = true
	SendCommand(CMD_REQ_REASON, true)
	if err = waitForResponce(); err != nil {
		return err
	}
	reason := string(SystemSt.ReasonBuf)
	log.Printf("Reason buf is %s\n", reason)
	if strings.HasPrefix(reason, "Button") {
		ProcButtonStart()
	} else if strings.HasPrefix(reason, "Sms") {
		ProcSetConfigStart()
	} else if strings.HasPrefix(reason, "Last") {
		ProcLastConfigStart()
	} else if strings.HasPrefix(reason, "Modem") {
		var cfg ModemPowerConfig
		// data := []byte(SystemSt.ReasonBuf)
		cfg.m1Pwr = SystemSt.ReasonBuf[6]
		cfg.m1Sim = SystemSt.ReasonBuf[7]
		cfg.m2Pwr = SystemSt.ReasonBuf[8]
		cfg.m2Sim = SystemSt.ReasonBuf[9]

		ProcModemStart(&cfg)
	}

	return nil
}

func procShutdown() {
	err := exec.Command("/bin/sh", "shutdown.sh").Run()
	if err != nil {
		log.Println(err)
	}
}

func procChangeOperator(idx uint8, operID string) error {

	flag.Set("address", router[idx].addr)
	flag.Set("username", router[idx].user)
	flag.Set("password", router[idx].pw)

	flag.Parse()

	c, err := routeros.Dial(*address, *username, *password)
	if err != nil {
		return err
	}

	reply, err := c.Run("/interface/lte/set", "=operator="+operID, "=.id=lte1")
	if err != nil {
		return err
	}

	// log.Printf("word %q\n", reply.Done.Word)
	if reply.Done.Word != "!done" {
		err = errors.New("Can't set operator")
		return err
	}

	time.Sleep(time.Second)

	reply, err = c.Run("/interface/lte/get", "=value-name=operator", "=number=lte1")
	if err != nil {
		return err
	}

	// log.Printf("key %q\n", reply.Done.List[0].Key)
	// log.Printf("val %q\n", reply.Done.List[0].Value)

	if reply.Done.List[0].Value != operID {
		err = errors.New("Can't verify operator")
		return err
	}

	c.Close()

	return nil
}

func modemTurnOn(idx uint8, sim uint8) error {
	var err error
	log.Printf("Modem turn on (num: %d, sim %d)\n", idx+1, sim)

	if PowerSt.Modem[idx] == true && ModemSt[idx].SimNum == sim {
		log.Println("Modem is working with this sim yet")
		return nil
	}

	if PowerSt.Modem[idx] == true {
		log.Println("\tPower off")
		FlagControlWaitResp = true
		SendObjectPwr(OBJECT_MODEM, idx, false)
		if err = waitForResponce(); err != nil {
			return err
		}
		time.Sleep(10 * time.Second)
	}

	log.Println("\tPower on")
	FlagControlWaitResp = true
	SendObjectPwr(OBJECT_MODEM, idx, true)
	if err = waitForResponce(); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	log.Println("\tFlightmode on")
	FlagControlWaitResp = true
	SendFlightmode(idx, true)
	if err = waitForResponce(); err != nil {
		return err
	}

	log.Println("\tChange sim")
	FlagControlWaitResp = true
	SendDoubleByte(CMD_CHANGE_SIM, idx, sim)
	if err = waitForResponce(); err != nil {
		return err
	}

	log.Println("\tLCD blink")
	FlagControlWaitResp = true
	SendDoubleByte(CMD_LCD_BLINK, idx, 0)
	if err = waitForResponce(); err != nil {
		return err
	}

	//? How to wait for modem loaded?
	time.Sleep(35 * time.Second)

	log.Println("\tReq modem info")
	FlagControlWaitResp = true
	SendShort(CMD_REQ_MODEM_INFO, idx)
	if err = waitForResponce(); err != nil {
		return err
	}
	//? Iccid should be changed on PCB by reading it from SIM?
	if modemStReq.Iccid != phFile.Bank[idx][sim-1].Iccid {
		log.Printf("\tIccid is wrong (recv: %s, file: %s)\n", modemStReq.Iccid, phFile.Bank[idx][sim-1].Iccid)
		err = errors.New("Iccid is wrong")
		return err
	}
	log.Printf("\tIccid is %s\n", modemStReq.Iccid)
	if modemStReq.Imei != phFile.Bank[idx][sim-1].Imei {
		time.Sleep(5 * time.Second)
		FlagControlWaitResp = true
		SendSetImei(idx, phFile.Bank[idx][sim-1].Imei)
		if err = waitForResponce(); err != nil {
			return err
		}

		FlagControlWaitResp = true
		SendShort(CMD_REQ_MODEM_INFO, idx)
		if err = waitForResponce(); err != nil {
			return err
		}

		if modemStReq.Imei != phFile.Bank[idx][sim-1].Imei {
			log.Printf("\tCan not set IMEI %s %s\n", modemStReq.Imei, phFile.Bank[idx][sim-1].Imei)
			err = errors.New("Can not set IMEI")
			return err
		}
	}
	log.Printf("\tIMEI is %s\n", modemStReq.Imei)
	time.Sleep(5 * time.Second)
	log.Printf("\tChanging operator to %s\n", phFile.Bank[idx][sim-1].OperID)

	err = procChangeOperator(idx, phFile.Bank[idx][sim-1].OperID)
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	time.Sleep(1 * time.Second)
	log.Println("\tFlightmode off")
	FlagControlWaitResp = true
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
		FlagControlWaitResp = true
		SendCommand(CMD_CTRL_ERROR, true)
		waitForResponce()
		FlagControlWaitResp = true
		SendCommand(CMD_PC_SHUTDOWN, true)
		waitForResponce()
		return err
	}

	SendConfig(cfg)

	return nil
}

// ProcModemStart - This procedure starts the modem
func ProcModemStart(cfg *ModemPowerConfig) {
	var err error

	reason := string(SystemSt.ReasonBuf)
	if !strings.HasPrefix(reason, "Last") {
		DeleteFile("config.txt")
	}
	SystemSt.ReasonBuf = nil

	FlagControlWaitResp = true
	SendShort(CMD_LOCK, 0)
	if err = waitForResponce(); err != nil {
		log.Printf("Cmd unlock: %q\n", err)
		return
	}

	if cfg.m2Pwr == 1 {
		err = modemTurnOn(1, cfg.m2Sim)
		if err != nil {
			log.Printf("Failed to turn on modem 2: %q\n", err)

			ModemSt[1].Iccid = ""
			ModemSt[1].Imei = ""
			ModemSt[1].Flymode = false
			ModemSt[1].SimNum = 0
			ModemSt[1].Phone = ""

			FlagControlWaitResp = true
			SendCommand(CMD_CFG_ERROR, true)
			waitForResponce()
		}
	}
	if err == nil {
		if cfg.m1Pwr == 1 {
			err = modemTurnOn(0, cfg.m1Sim)
			if err != nil {
				log.Printf("Failed to turn on modem 1: %q\n", err)

				ModemSt[0].Iccid = ""
				ModemSt[0].Imei = ""
				ModemSt[0].Flymode = false
				ModemSt[0].SimNum = 0
				ModemSt[0].Phone = ""

				FlagControlWaitResp = true
				SendCommand(CMD_CFG_ERROR, true)
				waitForResponce()
			}
		}
	}

	FlagControlWaitResp = true
	SendShort(CMD_UNLOCK, 0)
	if err = waitForResponce(); err != nil {
		log.Printf("Cmd unlock: %q\n", err)
		return
	}
}
