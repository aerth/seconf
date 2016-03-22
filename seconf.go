// The MIT License (MIT)
//
// Copyright (c) 2016 aerth
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// package seconf allows your software to store non-plaintext configuration files.
//
// check out the example application in _examples/hello for a working example.
//
//		 fmt.Println("Welcome to " + sn + ", " + configarray[0])
//		 fmt.Printf("Your %s is %s \n", os.Args[3], configarray[0])
//		 fmt.Printf("Your %s is %s \n", os.Args[4], configarray[1])
//		 fmt.Printf("Your %s is %s \n", os.Args[5], configarray[2])
package seconf

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"errors"

	"github.com/bgentry/speakeasy"
	"golang.org/x/crypto/nacl/secretbox"
)

const keySize = 32
const nonceSize = 24

// secustom is the filename that gets stored. for example, if secustom is "test", the configuration file will be saved as $HOME/.test
var secustom string
var username string
var password string
var hashbar = strings.Repeat("#", 80)

var configuser = ""
var configpass = ""

var configlock = ""

// Seconf is the struct for the seconf pathname and fields.
type Seconf struct {
	Id   int64
	Path string
	Args []string
}

/*
type Fielder struct {
	Id       int64
	Name     string
	Password bool
}
*/

// constainsString returns true if a slice contains a string.
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

// askForConfirmation returns true if the user types one of the "okayResponses"
// https://gist.github.com/albrow/5882501
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println(err)
		return false
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	quitResponses := []string{"q", "Q", "exit", "quit"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else if containsString(quitResponses, response) {
		return false
	} else {
		fmt.Println("\nNot valid answer, try again. [y/n] [yes/no]")
		return askForConfirmation()
	}
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// Prompt the user for the particular field.
func Prompt(header string) string {
	fmt.Printf("\n### " + header + " ###\n")
	fmt.Printf("\nPress ENTER when you are finished typing.\n\n")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		return line
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return Prompt(header)
	}
	return ""
}

// Create initializes a new configuration file, at $HOME/secustom with the title servicename and as many fields as needed. Any field starting with "pass" will be assumed a password and input will not be echoed.
func Create(secustom string, servicename string, arg ...string) {
	bar(servicename)
	configfields := &Seconf{
		Path: secustom,
		Args: arg,
	}

	var m1 map[int]string = map[int]string{}
	var newsplice []string
	for i := range configfields.Args {
		bar(servicename)
		if len(configfields.Args[i]) > 4 {
			if configfields.Args[i][0:4] == "pass" || configfields.Args[i][0:4] == "Pass" {
				//		fmt.Printf("\n### " + servicename + " ###\n")
				m1[i], _ = speakeasy.Ask(servicename + " " + configfields.Args[i] + ": ")
				if m1[i] == "" {
					bar(secustom)
					m1[i], _ = speakeasy.Ask(servicename + " " + configfields.Args[i] + ": ")
				}
				if m1[i] == "" {
					bar(secustom)
					m1[i], _ = speakeasy.Ask(servicename + " " + configfields.Args[i] + ": ")
				}
				if m1[i] == "" {
					bar(secustom)
					fmt.Println(configfields.Args[i] + " cannot be blank.")
					return
				}

			} else {
				m1[i] = Prompt(configfields.Args[i])
				if m1[i] == "" {
					bar(servicename)
					fmt.Println("Can not be blank.")
					m1[i] = Prompt(configfields.Args[i])
				}
				if m1[i] == "" {
					bar(servicename)
					fmt.Println("Can not be blank.")
					m1[i] = Prompt(configfields.Args[i])
				}
				if m1[i] == "" {
					bar(servicename)
					fmt.Println(configfields.Args[i] + " cannot be blank.")
					return
				}
			}
		} else {
			m1[i] = Prompt(configfields.Args[i])
		}
		newsplice = append(newsplice, m1[i]+"::::")
	}

	bar(servicename)
	configlock, _ := speakeasy.Ask("Create a password to encrypt config file:\nPress ENTER for no password.")
	var userKey = configlock
	var pad = []byte("«super jumpy fox jumps all over»")

	var messagebox = strings.Join(newsplice, "")
	messagebox = strings.TrimSuffix(messagebox, "::::")
	var message = []byte(messagebox)
	key := []byte(userKey)
	key = append(key, pad...)
	naclKey := new([keySize]byte)
	copy(naclKey[:], key[:keySize])
	nonce := new([nonceSize]byte)
	// Read bytes from random and put them in nonce until it is full.
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		fmt.Println("Could not read from random:", err)
		return
	}
	out := make([]byte, nonceSize)
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce, naclKey)
	err = ioutil.WriteFile(ReturnHome()+"/."+secustom, out, 0600)
	if err != nil {
		fmt.Println("Error while writing config file: ", err)
		return
	}
	fmt.Printf("Config file saved at "+ReturnHome()+"/."+secustom+" \nTotal size is %d bytes.\n",
		len(out))
	os.Exit(0)
}

// Detect returns TRUE if a seconf file exists.
func Detect(secustom string) bool {
	_, err := ioutil.ReadFile(ReturnHome() + "/." + secustom)
	if err != nil {
		return false
	}
	return true
}

// Read returns the decoded configuration file, or an error. Fields are separated by 4 colons. ("::::")
func Read(secustom string) (config string, err error) {
	bar(secustom)
	fmt.Println("Unlocking config file")
	configlock, err = speakeasy.Ask("Password: ")
	bar(secustom)
	var userKey = configlock
	var pad = []byte("«super jumpy fox jumps all over»")
	key := []byte(userKey)
	key = append(key, pad...)
	naclKey := new([keySize]byte)
	copy(naclKey[:], key[:keySize])
	nonce := new([nonceSize]byte)
	in, err := ioutil.ReadFile(ReturnHome() + "/." + secustom)
	if err != nil {
		fmt.Println(err)

	}
	copy(nonce[:], in[:nonceSize])
	configbytes, ok := secretbox.Open(nil, in[nonceSize:], nonce, naclKey)
	if !ok {
		fmt.Println("Could not decrypt the config file. Wrong password?")
		return
	}
	return string(configbytes), nil

}

// Cheap and effective way of clearing screen on unix. Ugly on windows.
func bar(secustom string) {
	versionbar := strings.Repeat("#", 10) + "\t" + secustom + "\t" + strings.Repeat("#", 30)
	print("\033[H\033[2J")
	fmt.Println(versionbar)
}

// ReturnHome is a cross-OS way of getting a HOMEDIR.
func ReturnHome() (homedir string) {
	homedir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if homedir == "" {
		homedir = os.Getenv("USERPROFILE")
	}
	if homedir == "" {
		homedir = os.Getenv("HOME")
	}
	return
}

// Lock() is the new version of Create(), It returns any errors during the process instead of using os.Exit()
func Lock(secustom string, servicename string, arg ...string) error {
	bar(servicename)
	configfields := &Seconf{
		Path: secustom,
		Args: arg,
	}

	var m1 map[int]string = map[int]string{}
	var newsplice []string
	for i := range configfields.Args {
		bar(servicename)
		if len(configfields.Args[i]) > 4 {
			if configfields.Args[i][0:4] == "pass" || configfields.Args[i][0:4] == "Pass" {
				m1[i], _ = speakeasy.Ask(servicename + " " + configfields.Args[i] + ": ")
				if m1[i] == "" {
					bar(servicename)
					m1[i], _ = speakeasy.Ask(servicename + " " + configfields.Args[i] + ": ")
				}
				if m1[i] == "" {
					bar(servicename)
					m1[i], _ = speakeasy.Ask(servicename + " " + configfields.Args[i] + ": ")
				}
				if m1[i] == "" {
					bar(servicename)
					return errors.New(configfields.Args[i] + " cannot be blank.")
				}

			} else {
				m1[i] = Prompt(configfields.Args[i])
				if m1[i] == "" {
					bar(servicename)
					fmt.Println("Can not be blank.")
					m1[i] = Prompt(configfields.Args[i])
				}
				if m1[i] == "" {
					bar(servicename)
					fmt.Println("Can not be blank.")
					m1[i] = Prompt(configfields.Args[i])
				}
				if m1[i] == "" {
					bar(servicename)
					return errors.New(configfields.Args[i] + " cannot be blank.")
				}
			}
		} else {
			m1[i] = Prompt(configfields.Args[i])
		}
		newsplice = append(newsplice, m1[i]+"::::")
	}
	bar(servicename)
	configlock, _ := speakeasy.Ask("Create a password to encrypt config file:\nPress ENTER for no password.")
	var userKey = configlock
	var pad = []byte("«super jumpy fox jumps all over»")

	var messagebox = strings.Join(newsplice, "")
	messagebox = strings.TrimSuffix(messagebox, "::::")
	var message = []byte(messagebox)
	key := []byte(userKey)
	key = append(key, pad...)
	naclKey := new([keySize]byte)
	copy(naclKey[:], key[:keySize])
	nonce := new([nonceSize]byte)
	// Read bytes from random and put them in nonce until it is full.
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return errors.New("Could not read from random: "+err.Error())
	}
	out := make([]byte, nonceSize)
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce, naclKey)
	err = ioutil.WriteFile(ReturnHome()+"/."+secustom, out, 0600)
	if err != nil {
		return errors.New("Error while writing config file: "+err.Error())
	}
	fmt.Printf("Config file saved at "+ReturnHome()+"/."+secustom+" \nTotal size is %d bytes.\n", len(out))
	return nil
}
