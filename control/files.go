package control

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var CfgFile FileConfig
var phFile FilePhones

func SetPhonesFile(str string) int {
	var record []string

	lines := strings.Split(str, ";")
	if len(lines) != 12 {
		return 1
	}
	for i := 0; i < len(lines); i++ {
		record = strings.Split(lines[i], ",")
		if len(record) != 3 {
			return 1
		}

		if i < 4 {
			phFile.Bank1[i].SimId = record[0]
			phFile.Bank1[i].Imei = record[1]
			phFile.Bank1[i].OperId = record[2]
		} else if i < 8 {
			phFile.Bank2[i-4].SimId = record[0]
			phFile.Bank2[i-4].Imei = record[1]
			phFile.Bank2[i-4].OperId = record[2]
		} else {
			phFile.Phones.PhonesOut[i-8] = record[0]
			phFile.Phones.PhonesIn[i-8] = record[1]
		}
	}

	return 0
}

func GetPhonesFile() string {
	var str string
	var record [3]string

	for i := 0; i < 12; i++ {
		if i < 4 {
			record[0] = phFile.Bank1[i].SimId
			record[1] = phFile.Bank1[i].Imei
			record[2] = phFile.Bank1[i].OperId
		} else if i < 8 {
			record[0] = phFile.Bank2[i-4].SimId
			record[1] = phFile.Bank2[i-4].Imei
			record[2] = phFile.Bank2[i-4].OperId
		} else {
			record[0] = phFile.Phones.PhonesOut[i-8]
			record[1] = phFile.Phones.PhonesIn[i-8]
			record[2] = ""
		}

		str += record[0] + "," + record[1] + "," + record[2] + "\n"
	}

	return str
}

func readPhonesFile(path string) (FilePhones, error) {
	var ph FilePhones

	csvfile, err := os.Open(path)
	if err != nil {
		return ph, err
	}
	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)

	// Iterate through the records
	for i := 0; i < 12; i++ {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return ph, err
		}

		if i < 4 {
			ph.Bank1[i].SimId = record[0]
			ph.Bank1[i].Imei = record[1]
			ph.Bank1[i].OperId = record[2]
		} else if i < 8 {
			ph.Bank2[i-4].SimId = record[0]
			ph.Bank2[i-4].Imei = record[1]
			ph.Bank2[i-4].OperId = record[2]
		} else {
			ph.Phones.PhonesOut[i-8] = record[0]
			ph.Phones.PhonesIn[i-8] = record[1]
		}
	}

	return ph, nil
}

func writePhonesFile(path string, ph FilePhones) error {
	var record [3]string
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	for i := 0; i < 12; i++ {
		if i < 4 {
			record[0] = ph.Bank1[i].SimId
			record[1] = ph.Bank1[i].Imei
			record[2] = ph.Bank1[i].OperId
		} else if i < 8 {
			record[0] = ph.Bank2[i-4].SimId
			record[1] = ph.Bank2[i-4].Imei
			record[2] = ph.Bank2[i-4].OperId
		} else {
			record[0] = ph.Phones.PhonesOut[i-8]
			record[1] = ph.Phones.PhonesIn[i-8]
			record[2] = ""
		}

		err := w.Write(record[:])
		if err != nil {
			return err
		}
	}

	return nil
}

func CfgToString(cfg FileConfig) string {
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

	for i := 0; i < 14; i++ {
		buf[i] = buf[i] + '0'
	}

	str := string(buf[:])
	return str
}

func StrToCfg(str string) FileConfig {
	var cfg FileConfig

	data := []byte(str)

	for i := 0; i < 14; i++ {
		data[i] = data[i] - '0'
	}

	if data[0] == 1 {
		cfg.Power.PowerStat = true
	}
	cfg.Power.BatLevel = 10*uint8(data[1]) + uint8(data[2])

	if data[3] == 1 {
		cfg.Power.Modem[0] = true
	}
	cfg.SimNum[0] = data[4]
	if data[5] == 1 {
		cfg.Power.Modem[1] = true
	}
	cfg.SimNum[1] = data[6]

	if data[7] == 1 {
		cfg.Power.Pc = true
	}
	if data[8] == 1 {
		cfg.Power.Wifi = true
	}
	if data[9] == 1 {
		cfg.Power.Relay[0] = true
	}
	if data[10] == 1 {
		cfg.Power.Relay[1] = true
	}

	if data[11] == 1 {
		cfg.ConfigErr = true
	}
	if data[12] == 1 {
		cfg.StateErr = true
	}
	if data[13] == 1 {
		cfg.ConnectErr = true
	}

	return cfg
}

func SetConfigFile(str string) {
	CfgFile = StrToCfg(str)
}

func GetConfigFile() string {
	return CfgToString(CfgFile)
}

func readConfigFile(path string) (FileConfig, error) {
	var cfg FileConfig
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	str := string(bs)
	fmt.Println("Read: ", str)

	cfg = StrToCfg(str)

	return cfg, nil
}

func writeConfigFile(path string, cfg FileConfig) error {
	str := CfgToString(cfg)
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
