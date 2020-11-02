package control

import (
	"testing"
)

func Test_cfg(t *testing.T) {
	cfg, err := readConfigFile("../config.txt")
	if err != nil {
		t.Fatal(err)
	}
	err = writeConfigFile("../config_wr.txt", cfg)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ph(t *testing.T) {
	ph, err := readPhonesFile("../phones.csv")
	if err != nil {
		t.Fatal(err)
	}
	err = writePhonesFile("../phones_wr.csv", ph)
	if err != nil {
		t.Fatal(err)
	}
}
