package control

import (
	"log"
	"time"

	"../com"
	"../crc16"
	sw "../go"
)

var table *crc16.Table

func Init() {
	log.Printf("Init protocol")

	com.Init(Callback)

	table = crc16.MakeTable(crc16.CRC16_MAXIM)

	com.Send([]byte("hellllo!\n"))
	//go comSend()
}

func SendCommand(cmdType byte, state bool) {
	var buf [6]byte

	buf[0] = cmdType
	buf[1] = 1
	if state {
		buf[2] = 1
	}

	crc := crc16.Checksum(buf[:], table)
	buf[3] = uint8(crc >> 8)
	buf[4] = uint8(crc & 0xff)
	buf[5] = byte('\n')

	com.Send(buf[:])
}

func SendData(cmdType byte, data []byte) {
	var dataLen = len(data)

	var buf = make([]byte, dataLen+5)

	buf[0] = cmdType
	for i := 0; i < dataLen; i++ {
		buf[2+i] = data[i]
	}

	crc := crc16.Checksum(buf[:], table)
	buf[3+dataLen] = uint8(crc >> 8)
	buf[4+dataLen] = uint8(crc & 0xff)
	buf[5+dataLen] = byte('\n')

	buf[1] = uint8(5 + dataLen)

	com.Send(buf[:])
}

func Callback(data []byte) {
	crc := crc16.Checksum(data, table)
	if crc != 0 {
		log.Printf("Bad crc16")
		// return
	}

	log.Printf("recieved: %q", data)

	switch data[0] {
	case CMD_LOCK:
		log.Printf("CMD_LOCK")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_UNLOCK:
		log.Printf("CMD_UNLOCK")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_FLYMODE:
		log.Printf("CMD_FLYMODE")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_POWER:
		log.Printf("CMD_POWER")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_CHANGE_SIM:
		log.Printf("CMD_CHANGE_SIM")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_LCD_PRINT:
		log.Printf("CMD_LCD_PRINT")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_SET_IMEI:
		log.Printf("CMD_SET_IMEI")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_SET_CONFIG:
		log.Printf("CMD_SET_CONFIG")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_CFG_ERROR:
		log.Printf("CMD_CFG_ERROR")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_CTRL_ERROR:
		log.Printf("CMD_CTRL_ERROR")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_PC_WAITMODE:
		log.Printf("CMD_PC_WAITMODE")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_PC_SHUTDOWN:
		log.Printf("CMD_PC_SHUTDOWN")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_PC_READY:
		log.Printf("CMD_PC_READY")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_NEW_PHONES:
		log.Printf("CMD_NEW_PHONES")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_REQ_MODEM_INFO:
		log.Printf("CMD_REQ_MODEM_INFO")

		var ptr int = 2
		idx := data[ptr]
		ptr++
		if data[ptr] == 1 {
			ModemSt[idx].Flymode = true
		} else {
			ModemSt[idx].Flymode = false
		}
		ptr++
		ModemSt[idx].SimNum = data[ptr]
		ptr++

		var simid = make([]byte, SIMID_SIZE)
		copy(simid, data[ptr:ptr+SIMID_SIZE])
		ModemSt[idx].SimId = string(simid)
		ptr += SIMID_SIZE

		var phone = make([]byte, PHONE_SIZE)
		copy(phone, data[ptr:ptr+PHONE_SIZE])
		ModemSt[idx].Phone = string(phone)
		ptr += PHONE_SIZE

		var imei = make([]byte, IMEI_SIZE)
		copy(imei, data[ptr:ptr+IMEI_SIZE])
		ModemSt[idx].Imei = string(imei)
		ptr += IMEI_SIZE

		ControlReqChan <- 1
	case CMD_REQ_PHONES:
		log.Printf("CMD_REQ_PHONES")

		var ptr int = 2
		ptr++
		for i := 0; i < 4; i++ {
			var phone = make([]byte, PHONE_SIZE)
			copy(phone, data[ptr:ptr+PHONE_SIZE])
			ModemPh.phonesOut[i] = string(phone)
			ptr += PHONE_SIZE
		}
		for i := 0; i < 4; i++ {
			var phone = make([]byte, PHONE_SIZE)
			copy(phone, data[ptr:ptr+PHONE_SIZE])
			ModemPh.phonesIn[i] = string(phone)
			ptr += PHONE_SIZE
		}
		ControlReqChan <- 1
	case CMD_REQ_REASON:
		log.Printf("CMD_REQ_REASON")

		len := data[1]
		copy(SystemSt.ReasonBuf, data[2:2+len])
		ControlReqChan <- 1
	case CMD_OUT_SHUTDOWN:
		log.Printf("CMD_OUT_SHUTDOWN")
	case CMD_OUT_SAVE_STATE:
		log.Printf("CMD_OUT_SAVE_STATE")
	case CMD_OUT_SIM_CHANGE:
		log.Printf("CMD_OUT_SIM_CHANGE")
	default:
	}
}

func comSend() {
	for i := 0; ; i++ {
		time.Sleep(5 * time.Second)

		com.Send([]byte("hellllooooo!\n"))
	}
}
