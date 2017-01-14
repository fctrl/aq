package aq

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Aq struct {
	ConfigDir string
	PinFile   string
}

func (a Aq) Reset() error {
	if a.ConfigDir == "" {
		return errors.New("reset without ConfigDir is not implemented")
	}
	if err := os.RemoveAll(a.ConfigDir); err != nil {
		return err
	}
	return os.RemoveAll(a.PinFile)
}

func (a Aq) AddUser(u User) error {
	var args []string
	if a.ConfigDir != "" {
		args = append(args, "-C", a.ConfigDir)
	}
	args = append(args, "adduser")
	if u.Name != "" {
		args = append(args, "-N", u.Name)
	}
	if u.ID != "" {
		args = append(args, "-u", u.ID)
	}
	if u.BankCode != "" {
		args = append(args, "-b", u.BankCode)
	}
	if u.ServerURL != "" {
		args = append(args, "-s", u.ServerURL)
	}
	if u.TokenType != "" {
		args = append(args, "-t", u.TokenType)
	}
	if u.HBCIVersion != "" {
		args = append(args, "--hbciversion="+u.HBCIVersion)
	}
	cmd := exec.Command("aqhbci-tool4", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, out)
	}
	return a.setHTTPVersion(u)
}

var httpRegexp = regexp.MustCompile("httpV(Major|Minor)=\"[\\d]\"")

func (a Aq) setHTTPVersion(u User) error {
	file, err := a.userSettingsFile(u)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	version := strings.Split(u.HTTPVersion, ".")
	if len(version) != 2 {
		return fmt.Errorf("invalid http version: %s", u.HTTPVersion)
	}
	matches := httpRegexp.FindAllSubmatch(data, -1)
	for _, m := range matches {
		var replace string
		kind := string(m[1])
		switch kind {
		case "Major":
			replace = version[0]
		case "Minor":
			replace = version[1]
		default:
			panic("bug")
		}
		replace = "httpV" + kind + "=\"" + replace + "\""
		data = bytes.Replace(data, m[0], []byte(replace), -1)
	}
	return ioutil.WriteFile(file, data, 0600)
}

func (a Aq) userSettingsFile(u User) (string, error) {
	userDir := filepath.Join(a.ConfigDir, "settings", "users", "*.conf")
	files, err := filepath.Glob(userDir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		if bytes.Index(data, []byte(u.ID)) == -1 {
			continue
		} else if bytes.Index(data, []byte(u.BankCode)) == -1 {
			continue
		}
		return file, nil
	}
	return "", errors.New("failed to locate user file")
}

func (a Aq) GetSysID(u User) error {
	if a.PinFile == "" {
		return errors.New("GetSysID without PinFile is not implemented")
	}
	pinData := fmt.Sprintf(`PIN_%s_%s = "%s"`, u.BankCode, u.ID, u.Pin)
	if err := ioutil.WriteFile(a.PinFile, []byte(pinData), 0600); err != nil {
		return err
	}
	var args []string
	if a.ConfigDir != "" { // TODO remove duplication
		args = append(args, "-C", a.ConfigDir)
	}
	args = append(args, "-P", a.PinFile)
	args = append(args, "-A", "-n")
	args = append(args, "getsysid")
	cmd := exec.Command("aqhbci-tool4", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, out)
	}
	return nil
}

type User struct {
	ID          string
	Name        string
	BankCode    string
	ServerURL   string
	TokenType   string
	HBCIVersion string
	HTTPVersion string
	Pin         string
}

type Account struct {
	ID       string
	BankCode string
}
