package control

import (
	"container/list"
	"log"
	"strings"
	"time"
	"unicode"

	"../com"
	"../crc16"
)

var table *crc16.Table

// InitProtocol - Init function
func InitProtocol() {
	log.Printf("Init protocol\n")

	com.Init(recieveHandler)

	SmsList = list.New()

	table = crc16.MakeMyTable(crc16.CRC16_MY)
}

func SendCommand(cmdType uint8, state bool) {
	// log.Printf("SendCommand\n")
	var buf [6]byte

	buf[0] = cmdType
	buf[1] = 1
	if state {
		buf[2] = 1
	}

	crc := crc16.Checksum(buf[:3], table)
	buf[3] = uint8(crc & 0xff)
	buf[4] = uint8(crc >> 8)
	buf[5] = byte('\n')

	com.Send(buf[:])
}

func SendShort(cmdType uint8, data byte) {
	// log.Printf("SendShort\n")
	var buf [6]byte

	buf[0] = cmdType
	buf[1] = 1
	buf[2] = data

	crc := crc16.Checksum(buf[:3], table)
	buf[3] = uint8(crc & 0xff)
	buf[4] = uint8(crc >> 8)
	buf[5] = byte('\n')

	com.Send(buf[:])
}

func SendData(cmdType uint8, data []byte) {
	// log.Printf("SendData\n")
	var dataLen = len(data)

	var buf = make([]byte, dataLen+5)

	buf[0] = cmdType
	buf[1] = uint8(dataLen)
	for i := 0; i < dataLen; i++ {
		buf[2+i] = data[i]
	}

	crc := crc16.Checksum(buf[0:len(buf)-3], table)
	buf[2+dataLen] = uint8(crc & 0xff)
	buf[3+dataLen] = uint8(crc >> 8)
	buf[4+dataLen] = byte('\n')

	com.Send(buf[:])
}

func SendDoubleByte(cmdType uint8, byte1 uint8, byte2 uint8) {
	var buf [2]byte

	buf[0] = byte1
	buf[1] = byte2

	SendData(cmdType, buf[:])
}

func SendFlightmode(idx uint8, state bool) {
	var buf [2]byte

	buf[0] = idx
	if state {
		buf[1] = 1
	}

	SendData(CMD_FLYMODE, buf[:])
}

func SendObjectPwr(obj uint8, idx uint8, state bool) {
	var buf [3]byte

	buf[0] = obj
	buf[1] = idx
	if state {
		buf[2] = 1
	}

	SendData(CMD_POWER, buf[:])
}

func SendSetImei(idx uint8, imei string) {

	len := 1 + IMEI_SIZE
	var buf = make([]byte, len)

	buf[0] = idx
	copy(buf[1:], imei)

	SendData(CMD_SET_IMEI, buf[:])
}

func SendConfig(cfg FileConfig) {
	var buf [14]byte

	if cfg.Power.PowerStat == true {
		buf[0] = 1
	}
	buf[1] = cfg.Power.BatLevel / 10
	buf[2] = cfg.Power.BatLevel % 10

	if cfg.Power.Modem[0] == true {
		buf[3] = 1
	}
	buf[4] = cfg.SimNum[0]

	if cfg.Power.Modem[1] == true {
		buf[5] = 1
	}
	buf[6] = cfg.SimNum[1]

	if cfg.Power.Pc == true {
		buf[7] = 1
	}
	if cfg.Power.Wifi == true {
		buf[8] = 1
	}

	if cfg.Power.Relay[0] == true {
		buf[9] = 1
	}
	if cfg.Power.Relay[1] == true {
		buf[10] = 1
	}
	if cfg.ConfigErr == true {
		buf[11] = 1
	}
	if cfg.StateErr == true {
		buf[12] = 1
	}
	if cfg.ConnectErr == true {
		buf[13] = 1
	}

	SendData(CMD_SET_CONFIG, buf[:])
}

func SendNewPhones(ph ModemPhones) {
	var buf [8 * PHONE_SIZE]byte

	var ptr int = 0
	for i := 0; i < 4; i++ {
		copy(buf[ptr:], ph.PhonesOut[i])
		ptr += PHONE_SIZE
	}
	for i := 0; i < 4; i++ {
		copy(buf[ptr:], ph.PhonesIn[i])
		ptr += PHONE_SIZE
	}

	SendData(CMD_NEW_PHONES, buf[:])
}

func SendSmsMessage(sms *SmsMessage) {
	len := 2 + PHONE_SIZE + len(sms.Message)
	var buf = make([]byte, len)

	var ptr int = 0

	// Modem num
	buf[ptr] = sms.ModemNum
	ptr++

	// Message type (now empty)
	buf[ptr] = sms.MsgType
	ptr++

	// Phone number
	copy(buf[ptr:], sms.Phone)
	ptr += PHONE_SIZE

	// Message
	copy(buf[ptr:], sms.Message)

	SendData(CMD_SEND_SMS, buf[:])
}

func recieveHandler(data []byte) {
	var crcIn uint16
	var crc [2]uint8
	if int(data[1]) != (len(data) - 5) {
		log.Printf("Wrong length %d (real %d)\n", data[1], (len(data) - 5))
		return
	}

	crcPkt := crc16.Checksum(data[:len(data)-3], table)

	crc[0] = uint8(crcPkt)
	crc[1] = uint8(crcPkt >> 8)
	if crc[0] == 0xFE {
		crc[0] = 0xFD
	}
	if crc[1] == 0xFE {
		crc[1] = 0xFD
	}
	crcPkt = uint16(crc[1]) << 8
	crcPkt += uint16(crc[0])

	crcIn = uint16(data[len(data)-2]) << 8
	crcIn += uint16(data[len(data)-3])

	if crcPkt != crcIn {
		log.Printf("Bad crc16 0x%X 0x%X\n", crcPkt, crcIn)
		return
	}
	// log.Printf("recv: ")
	// for i := 0; i < len(data)-1; i++ {
	// 	log.Printf("%02X ", data[i])
	// }
	// log.Printf("  \n")
	// // ! Return here bacause of there are blocking by channel below
	// return

	switch data[0] {
	case CMD_LOCK:
		// log.Printf("CMD_LOCK\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_UNLOCK:
		// log.Printf("CMD_UNLOCK\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_FLYMODE:
		// log.Printf("CMD_FLYMODE\n")
		var state bool
		idx := data[2]
		if data[3] == 0 {
			state = false
		} else {
			state = true
		}

		ModemSt[idx].Flymode = state

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- 1
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- 1
		}
	case CMD_POWER:
		// log.Printf("CMD_POWER\n")
		var state bool
		obj := data[2]
		idx := data[3]
		if data[4] == 0 {
			state = false
		} else {
			state = true
		}

		switch obj {
		case OBJECT_MODEM:
			PowerSt.Modem[idx] = state
		case OBJECT_WIFI:
			PowerSt.Wifi = state
		case OBJECT_RELAY:
			PowerSt.Relay[idx] = state
		}

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- 1
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- 1
		}
	case CMD_CHANGE_SIM:
		// log.Printf("CMD_CHANGE_SIM\n")
		idx := data[2]
		sim := data[3]
		var res uint8

		if sim > 0 {
			ModemSt[idx].SimNum = sim
			res = 1
		} else {
			res = 0
		}

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- res
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- res
		}
	case CMD_LCD_PRINT:
		// log.Printf("CMD_LCD_PRINT\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_LCD_BLINK:
		// log.Printf("CMD_LCD_BLINK\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_SET_IMEI:
		// log.Printf("CMD_SET_IMEI\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_SET_CONFIG:
		// log.Printf("CMD_SET_CONFIG\n")

		cfg, err := BytesToCfg(data[2 : 2+CONFIG_LEN])
		if err != nil {
			log.Printf("Error cfg recv %q\n", err)
			break
		}

		PowerSt.BatLevel = cfg.Power.BatLevel
		PowerSt.PowerStat = cfg.Power.PowerStat
		PowerSt.Pc = cfg.Power.Pc
		PowerSt.Wifi = cfg.Power.Wifi
		PowerSt.Relay[0] = cfg.Power.Relay[0]
		PowerSt.Relay[1] = cfg.Power.Relay[1]
		PowerSt.Modem[0] = cfg.Power.Modem[0]
		PowerSt.Modem[1] = cfg.Power.Modem[1]

		ModemSt[0].SimNum = cfg.SimNum[0]
		ModemSt[1].SimNum = cfg.SimNum[1]

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- 1
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- 1
		}
	case CMD_GET_CONFIG:
		// log.Printf("CMD_GET_CONFIG\n")

		cfg, err := BytesToCfg(data[2 : 2+CONFIG_LEN])
		if err != nil {
			log.Printf("Error cfg recv %q\n", err)
			break
		}

		PowerSt.BatLevel = cfg.Power.BatLevel
		PowerSt.PowerStat = cfg.Power.PowerStat
		PowerSt.Pc = cfg.Power.Pc
		PowerSt.Wifi = cfg.Power.Wifi
		PowerSt.Relay[0] = cfg.Power.Relay[0]
		PowerSt.Relay[1] = cfg.Power.Relay[1]
		PowerSt.Modem[0] = cfg.Power.Modem[0]
		PowerSt.Modem[1] = cfg.Power.Modem[1]

		ModemSt[0].SimNum = cfg.SimNum[0]
		ModemSt[1].SimNum = cfg.SimNum[1]

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- 1
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- 1
		}
	case CMD_CFG_ERROR:
		// log.Printf("CMD_CFG_ERROR\n")

		ModemSt[0].Iccid = ""
		ModemSt[0].Imei = ""
		ModemSt[0].Flymode = false
		ModemSt[0].SimNum = 0
		ModemSt[0].Phone = ""

		ModemSt[1].Iccid = ""
		ModemSt[1].Imei = ""
		ModemSt[1].Flymode = false
		ModemSt[1].SimNum = 0
		ModemSt[1].Phone = ""

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_CTRL_ERROR:
		// log.Printf("CMD_CTRL_ERROR\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_PC_WAITMODE:
		// log.Printf("CMD_PC_WAITMODE\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_PC_SHUTDOWN:
		// log.Printf("CMD_PC_SHUTDOWN\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_PC_READY:
		// log.Printf("CMD_PC_READY\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_NEW_PHONES:
		// log.Printf("CMD_NEW_PHONES\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_SEND_SMS:
		// log.Printf("CMD_SEND_SMS\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_REQ_MODEM_INFO:
		// log.Printf("CMD_REQ_MODEM_INFO\n")

		var st ModemStatus
		var ptr int = 2
		idx := data[ptr]
		ptr++
		if data[ptr] == 1 {
			st.Flymode = true
		} else {
			st.Flymode = false
		}
		ptr++
		st.SimNum = data[ptr]
		ptr++

		var iccid = make([]byte, ICCID_SIZE)
		copy(iccid, data[ptr:ptr+ICCID_SIZE])
		st.Iccid = string(iccid)
		st.Iccid = strings.Trim(st.Iccid, "\u0000")
		ptr += ICCID_SIZE

		var phone = make([]byte, PHONE_SIZE)
		copy(phone, data[ptr:ptr+PHONE_SIZE])
		st.Phone = string(phone)
		st.Phone = strings.Trim(st.Phone, "\u0000")
		ptr += PHONE_SIZE

		var imei = make([]byte, IMEI_SIZE)
		copy(imei, data[ptr:ptr+IMEI_SIZE])
		st.Imei = string(imei)
		st.Imei = strings.Trim(st.Imei, "\u0000")
		ptr += IMEI_SIZE

		if (idx >> 4) == 1 {
			idx = idx & 0x0F
			SmsModemSt[idx] = st
		} else {
			modemStReq = st
			ModemSt[idx] = st
		}

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- 1
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- 1
		}
	case CMD_REQ_CONN_INFO:
		// log.Printf("CMD_REQ_CONN_INFO\n")

		var st ModemConnStatus
		var ptr int = 2
		idx := data[ptr]
		ptr++

		st.Status = uint8(data[ptr])
		ptr++

		st.Csq = uint8(data[ptr])
		ptr++

		st.Tac = uint16(data[ptr+1]) << 8
		st.Tac += uint16(data[ptr])
		ptr += 2

		st.CellID = uint32(data[ptr+3]) << 24
		st.CellID += uint32(data[ptr+2]) << 16
		st.CellID += uint32(data[ptr+1]) << 8
		st.CellID += uint32(data[ptr])
		ptr += 4

		var oper = make([]byte, OPERID_SIZE)
		copy(oper, data[ptr:ptr+OPERID_SIZE])
		st.OperID = string(oper)
		st.OperID = strings.Trim(st.OperID, "\u0000")
		ptr += OPERID_SIZE

		if (idx >> 4) == 1 {
			idx = idx & 0x0F
			SmsConnSt[idx] = st
		} else {
			ConnSt[idx] = st
		}

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- 1
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- 1
		}
	case CMD_REQ_PHONES:
		// log.Printf("CMD_REQ_PHONES\n")

		var ph ModemPhones
		var ptr int = 2
		for i := 0; i < 4; i++ {
			var phone = make([]byte, PHONE_SIZE)
			copy(phone, data[ptr:ptr+PHONE_SIZE])
			ph.PhonesOut[i] = strings.TrimRightFunc(string(phone), func(r rune) bool {
				return !unicode.IsPunct(r) && !unicode.IsNumber(r)
			})
			ptr += PHONE_SIZE
		}
		for i := 0; i < 4; i++ {
			var phone = make([]byte, PHONE_SIZE)
			copy(phone, data[ptr:ptr+PHONE_SIZE])
			ph.PhonesIn[i] = strings.TrimRightFunc(string(phone), func(r rune) bool {
				return !unicode.IsPunct(r) && !unicode.IsNumber(r)
			})
			ptr += PHONE_SIZE
		}

		modemPhReq = ph
		if FlagControlWaitResp == true {
			ControlReqChan <- 1
		}
		if FlagHTTPWaitResp == true {
			HTTPReqChan <- 1
			FlagHTTPWaitResp = false
		}
	case CMD_REQ_REASON:
		// log.Printf("CMD_REQ_REASON\n")

		len := data[1]
		SystemSt.ReasonBuf = data[2 : 2+len]
		//copy(SystemSt.ReasonBuf, data[2:2+len])
		if FlagControlWaitResp == true {
			ControlReqChan <- 1
		}
	case CMD_OUT_SHUTDOWN:
		// log.Printf("CMD_OUT_SHUTDOWN\n")

		go procShutdown()
	case CMD_OUT_SAVE_STATE:
		// log.Printf("CMD_OUT_SAVE_STATE\n")

		cfg, err := BytesToCfg(data[2 : 2+CONFIG_LEN])
		if err != nil {
			log.Printf("Error cfg recv %q\n", err)
			return
		}
		CfgFile = cfg

		go writeConfigFile("config.txt", cfg)
	case CMD_OUT_SIM_CHANGE:
		// log.Printf("CMD_OUT_SIM_CHANGE\n")

		var cfg ModemPowerConfig
		cfg.m1Pwr = data[2]
		cfg.m1Sim = data[3]
		cfg.m2Pwr = data[4]
		cfg.m2Sim = data[5]

		go ProcModemStart(&cfg)
	case CMD_OUT_SMS:
		// log.Printf("CMD_OUT_SMS\n")

		var ptr uint8 = 2
		var sms SmsMessage

		sms.ModemNum = data[ptr]
		ptr++
		sms.MsgType = data[ptr]
		ptr++
		sms.Phone = string(data[ptr : ptr+PHONE_SIZE])
		sms.Phone = strings.Trim(sms.Phone, "\u0000")
		ptr = ptr + PHONE_SIZE
		msgLen := data[1] - PHONE_SIZE - 2
		sms.Message = string(data[ptr : ptr+msgLen])
		sms.Message = strings.Trim(sms.Message, "\r\n")

		//! sms may be cleared after end of function (make(sms, 1))
		if sms.MsgType == 1 {
			SmsList.PushBack(&sms)
		}
	default:
		log.Println("Unknown command")
	}
}

func comSend() {
	for i := 0; ; i++ {
		time.Sleep(5 * time.Second)

		com.Send([]byte("hellllooooo!\n"))
	}
}
