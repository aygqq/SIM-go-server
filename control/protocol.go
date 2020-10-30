package control

import (
	"fmt"
	"time"

	"../com"
	"../crc16"
	sw "../go"
)

var table *crc16.Table

func Init() {
	fmt.Printf("Init protocol\n")

	com.Init(Callback)

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
	buf[3] = uint8(crc >> 8)
	buf[4] = uint8(crc & 0xff)
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
	buf[3] = uint8(crc >> 8)
	buf[4] = uint8(crc & 0xff)
	buf[5] = byte('\n')

	com.Send(buf[:])
}

func SendData(cmdType uint8, data []byte) {
	fmt.Printf("SendData\n")
	var dataLen = len(data)

	var buf = make([]byte, dataLen+5)

	buf[0] = cmdType
	buf[1] = uint8(5 + dataLen)
	for i := 0; i < dataLen; i++ {
		buf[2+i] = data[i]
	}

	crc := crc16.Checksum(buf[0:len(buf)-3], table)
	buf[2+dataLen] = uint8(crc >> 8)
	buf[3+dataLen] = uint8(crc & 0xff)
	buf[4+dataLen] = byte('\n')

	com.Send(buf[:])
}

func SendDoubleByte(cmdType uint8, byte1 uint8, byte2 uint8) {
	var buf [2]byte

	buf[0] = byte1
	buf[1] = byte2

	SendData(CMD_CHANGE_SIM, buf[:])
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

func SendSimChange(bank uint8, sim uint8) {
	var buf [2]byte

	buf[0] = bank
	buf[1] = sim

	SendData(CMD_CHANGE_SIM, buf[:])
}

func SendLcdInfo(infoType uint8, info uint8) {
	var buf [2]byte

	buf[0] = infoType
	buf[1] = info

	SendData(CMD_LCD_PRINT, buf[:])
}

func SendLcdBlink(bank uint8, sim uint8) {
	var buf [2]byte

	buf[0] = bank
	buf[1] = sim

	SendData(CMD_LCD_BLINK, buf[:])
}

func SendSetImei(idx uint8, imei string) {
	var buf = make([]byte, 1+len(imei))

	buf[0] = idx
	copy(buf[1:], imei)

	SendData(CMD_SET_IMEI, buf[:])
}

func SendConfig(cfg FileConfig) {
	var buf [13]byte

	if cfg.power.PowerStat == true {
		buf[0] = 1
	}
	buf[1] = cfg.power.BatLevel

	if cfg.power.Pc == true {
		buf[2] = 1
	}
	if cfg.power.Wifi == true {
		buf[3] = 1
	}

	if cfg.power.Modem[0] == true {
		buf[4] = 1
	}
	buf[5] = cfg.simNum[0]

	if cfg.power.Modem[1] == true {
		buf[6] = 1
	}
	buf[7] = cfg.simNum[1]

	if cfg.power.Relay[0] == true {
		buf[8] = 1
	}
	if cfg.power.Relay[1] == true {
		buf[9] = 1
	}
	if cfg.configErr == true {
		buf[10] = 1
	}
	if cfg.stateErr == true {
		buf[11] = 1
	}
	if cfg.connectErr == true {
		buf[12] = 1
	}

	SendData(CMD_SET_CONFIG, buf[:])
}

func SendNewPhones(ph ModemPhones) {
	var buf [8 * PHONE_SIZE]byte

	var ptr int = 0
	for i := 0; i < 4; i++ {
		copy(buf[ptr:], ph.phonesOut[i])
		ptr += PHONE_SIZE
	}
	for i := 0; i < 4; i++ {
		copy(buf[ptr:], ph.phonesIn[i])
		ptr += PHONE_SIZE
	}

	SendData(CMD_NEW_PHONES, buf[:])
}

func Callback(data []byte) {
	crc := crc16.Checksum(data[:len(data)-1], table)

	if crc != 0 {
		fmt.Printf("Bad crc16 %X\n", crc)
		return
	}

	//! Return here bacause of there are blocking by channel below
	return

	switch data[0] {
	case CMD_LOCK:
		fmt.Printf("CMD_LOCK\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_UNLOCK:
		fmt.Printf("CMD_UNLOCK\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_FLYMODE:
		fmt.Printf("CMD_FLYMODE\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_POWER:
		fmt.Printf("CMD_POWER\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_CHANGE_SIM:
		fmt.Printf("CMD_CHANGE_SIM\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_LCD_PRINT:
		fmt.Printf("CMD_LCD_PRINT\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_LCD_BLINK:
		fmt.Printf("CMD_LCD_BLINK\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_SET_IMEI:
		fmt.Printf("CMD_SET_IMEI\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_SET_CONFIG:
		fmt.Printf("CMD_SET_CONFIG\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_CFG_ERROR:
		fmt.Printf("CMD_CFG_ERROR\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_CTRL_ERROR:
		fmt.Printf("CMD_CTRL_ERROR\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_PC_WAITMODE:
		fmt.Printf("CMD_PC_WAITMODE\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_PC_SHUTDOWN:
		fmt.Printf("CMD_PC_SHUTDOWN\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_PC_READY:
		fmt.Printf("CMD_PC_READY\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_NEW_PHONES:
		fmt.Printf("CMD_NEW_PHONES\n")

		if sw.FlagWaitResp == true {
			sw.HttpReqChan <- data[2]
			sw.FlagWaitResp = false
		} else {
			ControlReqChan <- data[2]
		}
	case CMD_REQ_MODEM_INFO:
		fmt.Printf("CMD_REQ_MODEM_INFO\n")

		var st [2]ModemStatus
		var ptr int = 2
		idx := data[ptr]
		ptr++
		if data[ptr] == 1 {
			st[idx].Flymode = true
		} else {
			st[idx].Flymode = false
		}
		ptr++
		st[idx].SimNum = data[ptr]
		ptr++

		var simid = make([]byte, SIMID_SIZE)
		copy(simid, data[ptr:ptr+SIMID_SIZE])
		st[idx].SimId = string(simid)
		ptr += SIMID_SIZE

		var phone = make([]byte, PHONE_SIZE)
		copy(phone, data[ptr:ptr+PHONE_SIZE])
		st[idx].Phone = string(phone)
		ptr += PHONE_SIZE

		var imei = make([]byte, IMEI_SIZE)
		copy(imei, data[ptr:ptr+IMEI_SIZE])
		st[idx].Imei = string(imei)
		ptr += IMEI_SIZE

		ModemSt[0] = st[0]
		ModemSt[1] = st[1]
		ControlReqChan <- 1
	case CMD_REQ_PHONES:
		fmt.Printf("CMD_REQ_PHONES\n")

		var ph ModemPhones
		var ptr int = 2
		for i := 0; i < 4; i++ {
			var phone = make([]byte, PHONE_SIZE)
			copy(phone, data[ptr:ptr+PHONE_SIZE])
			ph.phonesOut[i] = string(phone)
			ptr += PHONE_SIZE
		}
		for i := 0; i < 4; i++ {
			var phone = make([]byte, PHONE_SIZE)
			copy(phone, data[ptr:ptr+PHONE_SIZE])
			ph.phonesIn[i] = string(phone)
			ptr += PHONE_SIZE
		}

		ModemPh = ph
		ControlReqChan <- 1
	case CMD_REQ_REASON:
		fmt.Printf("CMD_REQ_REASON\n")

		len := data[1]
		copy(SystemSt.ReasonBuf, data[2:2+len])
		ControlReqChan <- 1
	case CMD_OUT_SHUTDOWN:
		fmt.Printf("CMD_OUT_SHUTDOWN\n")
		//TODO: Start algorithm
	case CMD_OUT_SAVE_STATE:
		fmt.Printf("CMD_OUT_SAVE_STATE\n")
		//TODO: Start algorithm
	case CMD_OUT_SIM_CHANGE:
		fmt.Printf("CMD_OUT_SIM_CHANGE\n")
		//TODO: Start algorithm
	default:
	}
}

func comSend() {
	for i := 0; ; i++ {
		time.Sleep(5 * time.Second)

		com.Send([]byte("hellllooooo!\n"))
	}
}
