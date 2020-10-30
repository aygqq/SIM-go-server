package control

import (
	"fmt"
	"io/ioutil"
	"os"
)

var cfgFile FileConfig

func readConfigFile(path string) (FileConfig, error) {
	var cfg FileConfig

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	str := string(bs)
	fmt.Println("Read: ", str)

	data := []byte(str)

	for i := 0; i < 14; i++ {
		data[i] = data[i] - '0'
	}

	if data[0] == 1 {
		cfg.power.PowerStat = true
	}
	cfg.power.BatLevel = 10*uint8(data[1]) + uint8(data[2])

	if data[3] == 1 {
		cfg.power.Modem[0] = true
	}
	cfg.simNum[0] = data[4]
	if data[5] == 1 {
		cfg.power.Modem[1] = true
	}
	cfg.simNum[1] = data[6]

	if data[7] == 1 {
		cfg.power.Pc = true
	}
	if data[8] == 1 {
		cfg.power.Wifi = true
	}
	if data[9] == 1 {
		cfg.power.Relay[0] = true
	}
	if data[10] == 1 {
		cfg.power.Relay[1] = true
	}

	if data[11] == 1 {
		cfg.configErr = true
	}
	if data[12] == 1 {
		cfg.stateErr = true
	}
	if data[13] == 1 {
		cfg.connectErr = true
	}

	return cfg, nil
}

func writeConfigFile(path string, cfg FileConfig) error {
	var buf [14]byte

	if cfg.power.PowerStat == true {
		buf[0] = 1
	}
	buf[1] = cfg.power.BatLevel / 10
	buf[2] = cfg.power.BatLevel % 10

	if cfg.power.Modem[0] == true {
		buf[3] = 1
	}
	buf[4] = cfg.simNum[0]

	if cfg.power.Modem[1] == true {
		buf[5] = 1
	}
	buf[6] = cfg.simNum[1]

	if cfg.power.Pc == true {
		buf[7] = 1
	}
	if cfg.power.Wifi == true {
		buf[8] = 1
	}

	if cfg.power.Relay[0] == true {
		buf[9] = 1
	}
	if cfg.power.Relay[1] == true {
		buf[10] = 1
	}
	if cfg.configErr == true {
		buf[11] = 1
	}
	if cfg.stateErr == true {
		buf[12] = 1
	}
	if cfg.connectErr == true {
		buf[13] = 1
	}

	for i := 0; i < 14; i++ {
		buf[i] = buf[i] + '0'
	}

	str := string(buf[:])
	fmt.Println("Write: ", str)

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Create error")
		return err
	}
	file.WriteString(str)
	file.Close()
	return nil
}
