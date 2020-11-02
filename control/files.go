package control

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var cfgFile FileConfig

func readPhonesFile(path string) (FilePhones, error) {
	var phFile FilePhones

	csvfile, err := os.Open(path)
	if err != nil {
		return phFile, err
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
			return phFile, err
		}

		if i < 4 {
			phFile.bank1[i].simId = record[0]
			phFile.bank1[i].imei = record[1]
			phFile.bank1[i].operId = record[2]
		} else if i < 8 {
			phFile.bank2[i-4].simId = record[0]
			phFile.bank2[i-4].imei = record[1]
			phFile.bank2[i-4].operId = record[2]
		} else {
			phFile.phones.phonesOut[i-8] = record[0]
			phFile.phones.phonesIn[i-8] = record[1]
		}
	}

	return phFile, nil
}

func writePhonesFile(path string, phFile FilePhones) error {
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
			record[0] = phFile.bank1[i].simId
			record[1] = phFile.bank1[i].imei
			record[2] = phFile.bank1[i].operId
		} else if i < 8 {
			record[0] = phFile.bank2[i-4].simId
			record[1] = phFile.bank2[i-4].imei
			record[2] = phFile.bank2[i-4].operId
		} else {
			record[0] = phFile.phones.phonesOut[i-8]
			record[1] = phFile.phones.phonesIn[i-8]
			record[2] = ""
		}

		err := w.Write(record[:])
		if err != nil {
			return err
		}
	}

	return nil
}

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
		cfg.power.powerStat = true
	}
	cfg.power.batLevel = 10*uint8(data[1]) + uint8(data[2])

	if data[3] == 1 {
		cfg.power.modem[0] = true
	}
	cfg.simNum[0] = data[4]
	if data[5] == 1 {
		cfg.power.modem[1] = true
	}
	cfg.simNum[1] = data[6]

	if data[7] == 1 {
		cfg.power.pc = true
	}
	if data[8] == 1 {
		cfg.power.wifi = true
	}
	if data[9] == 1 {
		cfg.power.relay[0] = true
	}
	if data[10] == 1 {
		cfg.power.relay[1] = true
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

	if cfg.power.powerStat == true {
		buf[0] = 1
	}
	buf[1] = cfg.power.batLevel / 10
	buf[2] = cfg.power.batLevel % 10

	if cfg.power.modem[0] == true {
		buf[3] = 1
	}
	buf[4] = cfg.simNum[0]

	if cfg.power.modem[1] == true {
		buf[5] = 1
	}
	buf[6] = cfg.simNum[1]

	if cfg.power.pc == true {
		buf[7] = 1
	}
	if cfg.power.wifi == true {
		buf[8] = 1
	}

	if cfg.power.relay[0] == true {
		buf[9] = 1
	}
	if cfg.power.relay[1] == true {
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
