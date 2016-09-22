// Package seconf allows your software to store non-plaintext configuration files.
package seconf

/*
The MIT License (MIT)

Copyright (c) 2016 aerth

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bgentry/speakeasy"
	"golang.org/x/crypto/nacl/secretbox"
)

func init() {
	// Hopefully a clean exit
	interrupt := make(chan os.Signal, 1)
	stop := make(chan os.Signal, 1)
	reload := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(stop, syscall.SIGINT)
	signal.Notify(reload, syscall.SIGHUP)

	go func() {
		select {
		case signal := <-interrupt:
			fmt.Println("Got signal:", signal)
			fmt.Println("Dying")
			os.Exit(0)
		case signal := <-reload:
			fmt.Println("Got signal:", signal)
			fmt.Println("Dying")
			os.Exit(0)
		case signal := <-stop:
			fmt.Printf("Got signal:%v\n", signal)
			fmt.Println("Dying")
			os.Exit(0)
		}
	}()
}
func encrypt(filename string, configger Config) error {
	configlock, _ := speakeasy.Ask("Create a password to encrypt config file:\nPress ENTER for no password.")
	var userKey = configlock

	message, err := json.Marshal(configger)
	key := []byte(userKey)
	key = append(key, pad...)
	naclKey := new([keySize]byte)
	copy(naclKey[:], key[:keySize])
	nonce := new([nonceSize]byte)
	// Read bytes from random and put them in nonce until it is full.
	_, err = io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return errors.New("Could not read from random: " + err.Error())
	}
	out := make([]byte, nonceSize)
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce, naclKey)
	err = ioutil.WriteFile(filename, out, 0600)
	if err != nil {
		return err
	}
	return nil
}
func decrypt(filename string) ([]byte, error) {

	naclKey := new([keySize]byte)
	copy(naclKey[:], pad[:keySize])
	nonce := new([nonceSize]byte)
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	copy(nonce[:], in[:nonceSize])
	configbytes, ok := secretbox.Open(nil, in[nonceSize:], nonce, naclKey)
	if ok {
		return configbytes, nil
	}
	return nil, errors.New("Can't decrypt.")
}

// Read returns the decoded configuration file, or an error. Fields are separated by 4 colons. ("::::")
func ReadJSON(secustom string) (configger Config, err error) {
	// This is the default encoded-but-not-safe blank password

	naclKey := new([keySize]byte)
	copy(naclKey[:], pad[:keySize])
	nonce := new([nonceSize]byte)
	in, err := ioutil.ReadFile(secustom)
	if err != nil {
		return configger, err
	}
	copy(nonce[:], in[:nonceSize])
	configbytes, ok := secretbox.Open(nil, in[nonceSize:], nonce, naclKey)
	if ok {

		err = json.Unmarshal(configbytes, &configger)
		if err != nil {
			return configger, err
		}
		return configger, nil

	}

	// The blank password didn't work. Ask the user via speakeasy
	configlock, err = speakeasy.Ask("Password: ")
	var userKey = configlock
	key := []byte(userKey)
	key = append(key, pad...)
	naclKey = new([keySize]byte)
	copy(naclKey[:], key[:keySize])
	nonce = new([nonceSize]byte)
	in, err = ioutil.ReadFile(secustom)
	if err != nil {
		return configger, err
	}
	copy(nonce[:], in[:nonceSize])
	configbytes, ok = secretbox.Open(nil, in[nonceSize:], nonce, naclKey)
	if !ok {
		return configger, errors.New("Could not decrypt the config file. Wrong password?")
	}

	err = json.Unmarshal(configbytes, &configger)
	if err != nil {
		return configger, err
	}
	return configger, nil

}

type Config struct {
	Fields map[string]interface{} `json:",string,omitempty"`
}

// Lock() is the new version of Create(), It returns any errors during the process instead of using os.Exit()
func LockJSON(secustom string, servicename string, field map[string]string) error {
	if field == nil {
		return errors.New("Fields cant be nil")
	}

	if secustom == "" {
		return errors.New("Seconf location cant be blank")
	}

	if servicename == "" {
		servicename = "config"
	}

	servicename = "[" + servicename + "]"
	var m1 = map[string]interface{}{}

	for i, k := range field {
		if strings.Contains(i, "pass") || strings.Contains(i, "Pass") || strings.Contains(i, "Key") || strings.Contains(i, "key") || i[0:4] == "pass" || i[0:4] == "Pass" {
			// Is a password field
			m1[i], _ = speakeasy.Ask(servicename + " " + k + ": ")
			if m1[i] == "" {
				m1[i], _ = speakeasy.Ask(servicename + " " + k + ": ") // blank 1
			}
			if m1[i] == "" {
				m1[i], _ = speakeasy.Ask(servicename + " " + k + ": ") // blank 2
			}
			if m1[i] == "" {
				return errors.New(i + " cannot be blank.") // User gave blank
			}

		} else {

			m1[i] = Prompt(servicename + " " + k)
			if m1[i] == "" {

				fmt.Println("Can not be blank.")
				m1[i] = Prompt(servicename + " " + k)
			}
			if m1[i] == "" {

				fmt.Println("Can not be blank.")
				m1[i] = Prompt(servicename + " " + k)
			}
			if m1[i] == "" {

				return errors.New(i + " cannot be blank.")
			}
		}

		//newsplice = append(newsplice, m1[k].(string)+"::::"
	}

	configlock, _ := speakeasy.Ask("\nCreate a password to encrypt config file:\nPress ENTER for no password.\n")
	var userKey = configlock

	var configger Config
	configger.Fields = m1
	message, err := json.Marshal(configger)
	key := []byte(userKey)
	key = append(key, pad...)
	naclKey := new([keySize]byte)
	copy(naclKey[:], key[:keySize])
	nonce := new([nonceSize]byte)
	// Read bytes from random and put them in nonce until it is full.
	_, err = io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return errors.New("Could not read from random: " + err.Error())
	}
	out := make([]byte, nonceSize)
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce, naclKey)
	err = ioutil.WriteFile(secustom, out, 0600)
	if err != nil {
		return errors.New("Error while writing config file: " + err.Error())
	}
	fmt.Printf("Config file v2 saved at %q\nTotal size is %d bytes.\n", secustom, len(out))
	return nil
}

// Exists returns TRUE if a seconf file exists. (absolute or relative path)
func Exists(secustom string) bool {

	_, err := ioutil.ReadFile(secustom)
	if err != nil {
		return false
	}
	return true
}

func Destroy(secustom string) error {
	if secustom == "" {
		return errors.New("No filename")
	}
	if !Detect(secustom) {
		return errors.New("File not found")
	}
	err := os.Remove(returnHome() + "/." + secustom)
	if err != nil {
		return err
	}
	return nil
}
