package control

import (
	"container/list"
	"fmt"
	"strings"
	"time"
	"unicode"

	"../com"
	"../crc16"
)

var table *crc16.Table

// InitProtocol - Init function
func InitProtocol() {
	fmt.Printf("Init protocol\n")

	com.Init(recieveHandler)

	SmsList = list.New()

	//! TODO: Table must be simmilar with PCB's table
	table = crc16.MakeMyTable(crc16.CRC16_MY)
}

func SendCommand(cmdType uint8, state bool) {
	fmt.Printf("SendCommand\n")
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
	fmt.Printf("SendShort\n")
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
	fmt.Printf("SendData\n")
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

	SendData(CMD_NEW_PHONES, buf[:])
}

func recieveHandler(data []byte) {
	if int(data[1]) != (len(data) - 5) {
		fmt.Printf("Wrong length %d (real %d)\n", data[1], (len(data) - 4))
		return
	}

	crc := crc16.Checksum(data[:len(data)-1], table)

	if crc != 0 {
		fmt.Printf("Bad crc16 %X\n", crc)
		return
	}
	// fmt.Printf("recv: ")
	// for i := 0; i < len(data)-1; i++ {
	// 	fmt.Printf("%02X ", data[i])
	// }
	// fmt.Printf("  \n")
	// // ! Return here bacause of there are blocking by channel below
	// return

	switch data[0] {
	case CMD_LOCK:
		// fmt.Printf("CMD_LOCK\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_UNLOCK:
		// fmt.Printf("CMD_UNLOCK\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_FLYMODE:
		// fmt.Printf("CMD_FLYMODE\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_POWER:
		// fmt.Printf("CMD_POWER\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_CHANGE_SIM:
		// fmt.Printf("CMD_CHANGE_SIM\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_LCD_PRINT:
		// fmt.Printf("CMD_LCD_PRINT\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_LCD_BLINK:
		// fmt.Printf("CMD_LCD_BLINK\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_SET_IMEI:
		// fmt.Printf("CMD_SET_IMEI\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_SET_CONFIG:
		// fmt.Printf("CMD_SET_CONFIG\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_CFG_ERROR:
		// fmt.Printf("CMD_CFG_ERROR\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_CTRL_ERROR:
		// fmt.Printf("CMD_CTRL_ERROR\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_PC_WAITMODE:
		// fmt.Printf("CMD_PC_WAITMODE\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_PC_SHUTDOWN:
		// fmt.Printf("CMD_PC_SHUTDOWN\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_PC_READY:
		// fmt.Printf("CMD_PC_READY\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_NEW_PHONES:
		// fmt.Printf("CMD_NEW_PHONES\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_SEND_SMS:
		// fmt.Printf("CMD_SEND_SMS\n")

		if FlagHTTPWaitResp == true {
			HTTPReqChan <- data[2]
			FlagHTTPWaitResp = false
		}
		if FlagControlWaitResp == true {
			ControlReqChan <- data[2]
		}
	case CMD_REQ_MODEM_INFO:
		// fmt.Printf("CMD_REQ_MODEM_INFO\n")

		var st ModemStatus
		var ptr int = 2
		//idx := data[ptr]
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
		ptr += ICCID_SIZE

		var phone = make([]byte, PHONE_SIZE)
		copy(phone, data[ptr:ptr+PHONE_SIZE])
		st.Phone = string(phone)
		ptr += PHONE_SIZE

		var imei = make([]byte, IMEI_SIZE)
		copy(imei, data[ptr:ptr+IMEI_SIZE])
		st.Imei = string(imei)
		ptr += IMEI_SIZE

		modemStReq = st
		if FlagControlWaitResp == true {
			ControlReqChan <- 1
		}
	case CMD_REQ_PHONES:
		// fmt.Printf("CMD_REQ_PHONES\n")

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
	case CMD_REQ_REASON:
		// fmt.Printf("CMD_REQ_REASON\n")

		len := data[1]
		SystemSt.ReasonBuf = string(data[2 : 2+len])
		//copy(SystemSt.ReasonBuf, data[2:2+len])
		if FlagControlWaitResp == true {
			ControlReqChan <- 1
		}
	case CMD_OUT_SHUTDOWN:
		// fmt.Printf("CMD_OUT_SHUTDOWN\n")

		go procShutdown()
	case CMD_OUT_SAVE_STATE:
		// fmt.Printf("CMD_OUT_SAVE_STATE\n")

		str := string(data[2 : 2+CONFIG_LEN])
		cfg, err := StrToCfg(str)
		if err != nil {
			return
		}
		CfgFile = cfg

		go writeConfigFile("../config.txt", cfg)
	case CMD_OUT_SIM_CHANGE:
		// fmt.Printf("CMD_OUT_SIM_CHANGE\n")

		var cfg ModemPowerConfig
		cfg.m1Pwr = data[2]
		cfg.m1Sim = data[3]
		cfg.m2Pwr = data[4]
		cfg.m2Sim = data[5]

		go ProcModemStart(&cfg)
	case CMD_OUT_SMS:
		// fmt.Printf("CMD_OUT_SMS\n")

		var ptr uint8 = 2
		var sms SmsMessage

		sms.ModemNum = data[ptr]
		ptr++
		sms.MsgType = data[ptr]
		ptr++
		sms.Phone = string(data[ptr : ptr+PHONE_SIZE])
		ptr = ptr + PHONE_SIZE
		msgLen := data[1] - PHONE_SIZE - 2
		sms.Message = string(data[ptr : ptr+msgLen])

		//! sms may be cleared after end of function (make(sms, 1))
		if sms.MsgType == 1 {
			SmsList.PushBack(&sms)
		}
	default:
		fmt.Println("Unknown command")
	}
}

func comSend() {
	for i := 0; ; i++ {
		time.Sleep(5 * time.Second)

		com.Send([]byte("hellllooooo!\n"))
	}
}
