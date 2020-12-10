package control

import (
	"encoding/csv"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

var CfgFile FileConfig
var phFile FilePhones

func checkIccid(str string) error {
	data := []byte(str)
	err := errors.New("Failed to parse Iccid")

	if len(data) != ICCID_SIZE {
		return err
	}

	for i := 0; i < ICCID_SIZE; i++ {
		if data[i] < '0' || data[i] > '9' {
			return err
		}
	}

	return nil
}

func checkImei(str string) error {
	data := []byte(str)
	err := errors.New("Failed to parse IMEI")

	if len(data) != IMEI_SIZE {
		return err
	}

	for i := 0; i < IMEI_SIZE; i++ {
		if data[i] < '0' || data[i] > '9' {
			return err
		}
	}

	return nil
}

func checkOperID(str string) error {
	data := []byte(str)
	err := errors.New("Failed to parse OperID")

	if len(data) != OPERID_SIZE {
		return err
	}

	for i := 0; i < OPERID_SIZE; i++ {
		if data[i] < '0' || data[i] > '9' {
			return err
		}
	}

	return nil
}

func checkPhone(str string) error {
	data := []byte(str)
	err := errors.New("Failed to parse IMEI")
	isEmpty := true
	isWrong := false

	if len(data) > PHONE_SIZE {
		return err
	}

	for i := 0; i < len(data); i++ {
		if data[i] != '*' {
			isEmpty = false
		}
	}

	for i := 0; i < len(data); i++ {
		if (data[i] < '0' || data[i] > '9') && data[i] != '+' {
			isWrong = true
		}
	}

	if isWrong == true && isEmpty == false {
		return err
	}

	return nil
}

func checkPhonesFile(file *FilePhones) error {
	var err error

	for i := 0; i < 4; i++ {
		err = checkIccid(file.Bank[0][i].Iccid)
		if err != nil {
			return err
		}
		err = checkImei(file.Bank[0][i].Imei)
		if err != nil {
			return err
		}
		err = checkOperID(file.Bank[0][i].OperID)
		if err != nil {
			return err
		}
		err = checkIccid(file.Bank[1][i].Iccid)
		if err != nil {
			return err
		}
		err = checkImei(file.Bank[1][i].Imei)
		if err != nil {
			return err
		}
		err = checkOperID(file.Bank[0][i].OperID)
		if err != nil {
			return err
		}
		err = checkPhone(file.Phones.PhonesIn[i])
		if err != nil {
			return err
		}
		err = checkPhone(file.Phones.PhonesOut[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func SetPhonesFile(records *[12][3]string) int {
	log.Println("SetPhonesFile")
	for i := 0; i < 12; i++ {
		if i < 4 {
			phFile.Bank[0][i].Iccid = records[i][0]
			phFile.Bank[0][i].Imei = records[i][1]
			phFile.Bank[0][i].OperID = records[i][2]
		} else if i < 8 {
			phFile.Bank[1][i-4].Iccid = records[i][0]
			phFile.Bank[1][i-4].Imei = records[i][1]
			phFile.Bank[1][i-4].OperID = records[i][2]
		} else {
			phFile.Phones.PhonesOut[i-8] = records[i][0]
			phFile.Phones.PhonesIn[i-8] = records[i][1]
		}
	}

	return 0
}

func GetPhonesFile(records *[12][3]string) {
	log.Println("GetPhonesFile")
	for i := 0; i < 12; i++ {
		if i < 4 {
			records[i][0] = phFile.Bank[0][i].Iccid
			records[i][1] = phFile.Bank[0][i].Imei
			records[i][2] = phFile.Bank[0][i].OperID
		} else if i < 8 {
			records[i][0] = phFile.Bank[1][i-4].Iccid
			records[i][1] = phFile.Bank[1][i-4].Imei
			records[i][2] = phFile.Bank[1][i-4].OperID
		} else {
			records[i][0] = phFile.Phones.PhonesOut[i-8]
			records[i][1] = phFile.Phones.PhonesIn[i-8]
			records[i][2] = ""
		}
	}
}

func readPhonesFile(path string) (FilePhones, error) {
	log.Println("readPhonesFile")
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
		if err != nil {
			return ph, err
		}

		if i < 4 {
			ph.Bank[0][i].Iccid = record[0]
			ph.Bank[0][i].Imei = record[1]
			ph.Bank[0][i].OperID = record[2]
		} else if i < 8 {
			ph.Bank[1][i-4].Iccid = record[0]
			ph.Bank[1][i-4].Imei = record[1]
			ph.Bank[1][i-4].OperID = record[2]
		} else {
			ph.Phones.PhonesOut[i-8] = record[0]
			ph.Phones.PhonesIn[i-8] = record[1]
		}
	}

	err = checkPhonesFile(&ph)

	return ph, err
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
			record[0] = ph.Bank[0][i].Iccid
			record[1] = ph.Bank[0][i].Imei
			record[2] = ph.Bank[0][i].OperID
		} else if i < 8 {
			record[0] = ph.Bank[1][i-4].Iccid
			record[1] = ph.Bank[1][i-4].Imei
			record[2] = ph.Bank[1][i-4].OperID
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

func StrToCfg(str string) (FileConfig, error) {
	var cfg FileConfig

	data := []byte(str)

	if len(data) != CONFIG_LEN {
		return cfg, errors.New("Config file format error")
	}

	for i := 0; i < CONFIG_LEN; i++ {
		if data[i] < '0' || data[i] > '9' {
			return cfg, errors.New("Config file format error")
		}
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

	return cfg, nil
}

func BytesToCfg(data []byte) (FileConfig, error) {
	var cfg FileConfig

	if len(data) != CONFIG_LEN {
		return cfg, errors.New("Config file format error")
	}

	for i := 0; i < CONFIG_LEN; i++ {
		if data[i] < 0 || data[i] > 9 {
			return cfg, errors.New("Config format error")
		}
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

	return cfg, nil
}

func SetConfigFile(str string) error {
	cfg, err := StrToCfg(str)
	CfgFile = cfg
	return err
}

func GetConfigFileString() string {
	return CfgToString(CfgFile)
}

func GetConfigFile() FileConfig {
	return CfgFile
}

func readConfigFile(path string) (FileConfig, error) {
	var cfg FileConfig
	var err error

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	str := string(bs)
	log.Println("Read: ", str)

	cfg, err = StrToCfg(str)

	return cfg, err
}

func writeConfigFile(path string, cfg FileConfig) error {
	str := CfgToString(cfg)
	log.Println("Write: ", str)

	file, err := os.Create(path)
	if err != nil {
		log.Printf("Create error")
		return err
	}
	file.WriteString(str)
	file.Close()
	return nil
}
