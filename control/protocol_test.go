package control

import (
	"fmt"
	"testing"
	"time"

	"../com"
	"../crc16"
)

var flag bool = false

func Test(t *testing.T) {
	fmt.Printf("Init protocol\n")

	com.Init(Callback_test)

	//! TODO: Table must be simmilar with PCB's table
	table = crc16.MakeMyTable(crc16.CRC16_MY)

	flag = false
	SendCommand(CMD_PC_READY, true)
	time.Sleep(time.Second)
	if flag == false {
		t.Fatal("Send command failed")
	}

	flag = false
	SendShort(CMD_LOCK, 2)
	time.Sleep(time.Second)
	if flag == false {
		t.Fatal("Send short failed")
	}

	flag = false
	SendSetImei(1, "1234567812345678")
	time.Sleep(time.Second)
	if flag == false {
		t.Fatal("Send imei failed")
	}

	flag = false
	var cfg FileConfig
	cfg.ConnectErr = true
	cfg.Power.Pc = true
	cfg.Power.Modem[0] = true
	SendConfig(cfg)
	time.Sleep(time.Second)
	if flag == false {
		t.Fatal("Send config failed")
	}

	flag = false
	var ph ModemPhones
	ph.PhonesOut[0] = "111111111111"
	ph.PhonesOut[1] = "222222222222"
	ph.PhonesOut[2] = "333333333333"
	ph.PhonesIn[0] = "+111111111111"
	ph.PhonesIn[1] = "+222222222222"
	ph.PhonesIn[2] = "+333333333333"
	ph.PhonesIn[3] = "99999999999999999999999999"
	SendNewPhones(ph)
	time.Sleep(time.Second)
	if flag == false {
		t.Fatal("Send phones failed")
	}
}

func Callback_test(data []byte) {
	crc := crc16.Checksum(data[:len(data)-1], table)

	var crcIn uint16
	crcIn = uint16(data[len(data)-3]) << 8
	crcIn += uint16(data[len(data)-2])

	fmt.Printf("recv: ")
	for i := 0; i < len(data)-1; i++ {
		fmt.Printf("%02X ", data[i])
	}
	fmt.Printf("  \n")

	if crc != 0 {
		fmt.Printf("Bad crc16 %X %X\n", crc, crcIn)
		return
	}
	fmt.Printf("Good crc16 %X %X\n", crc, crcIn)
	flag = true
	return
}
