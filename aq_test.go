package aq

import (
	"os"
	"testing"
)

func TestAddUser(t *testing.T) {
	a := Aq{ConfigDir: "./aqconf", PinFile: "./aqpin"}
	if err := a.Reset(); err != nil {
		t.Fatal(err)
	}
	u := User{
		Name:        os.Getenv("AQ_USER_NAME"),
		ID:          os.Getenv("AQ_USER_ID"),
		BankCode:    os.Getenv("AQ_BANK_CODE"),
		ServerURL:   os.Getenv("AQ_SERVER_URL"),
		TokenType:   os.Getenv("AQ_TOKEN_TYPE"),
		HBCIVersion: os.Getenv("AQ_HBCI_VERSION"),
		HTTPVersion: os.Getenv("AQ_HTTP_VERSION"),
		Pin:         os.Getenv("AQ_PIN"),
	}
	if err := a.AddUser(u); err != nil {
		t.Fatal(err)
	} else if err := a.GetSysID(u); err != nil {
		t.Fatal(err)
	}
}
