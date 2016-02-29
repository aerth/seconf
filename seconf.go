package seconf

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"github.com/bgentry/speakeasy"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const keySize = 32
const nonceSize = 24

var secustom string
var username = os.Getenv("USER")
var password = os.Getenv("SECONFPASS")
var hashbar = strings.Repeat("#", 80)

var configuser = ""
var configpass = ""

var configlock = ""

type Seconf struct {
	Id   int64
	Path string
	Args []string
}
type Fielder struct {
	Id       int64
	Name     string
	Password bool
}

/*
func main() {
// command: seconf create
if os.Args[1] == "config" {

  if Detect(s) == false {
    fmt.Println("Creating config file. You will be asked for your user,and password.")
    fmt.Println("Your password will NOT echo.\n")
    Create()
  } else {
    fmt.Println("Config file already exists.\nIf you want to create a new config file, move or delete the existing one.\n")
    fmt.Println(os.Getenv("HOME") + "/."+secustom+"\n")
    os.Exit(1)
  }
}

}

*/
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func Prompt(header string) string {
	fmt.Printf("\n### " + header + " ###\n")
	fmt.Printf("\nPress ENTER when you are finished typing.\n\n")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		//	fmt.Println(line)
		return line
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return ""
}

func Create(secustom string, servicename string, arg ...string) {
	bar(secustom)
	configfields := &Seconf{
		Path: secustom,
		Args: arg,
	}

	var m1 map[int]string = map[int]string{}
	var newsplice []string
	for i := range configfields.Args {
		bar(secustom)
		if len(configfields.Args[i]) > 4 {
			if configfields.Args[i][0:4] == "pass" {
				fmt.Printf("\n### " + servicename + " ###\n")
				m1[i], _ = speakeasy.Ask(servicename + " " + configfields.Args[i] + ":")
			} else {
				m1[i] = Prompt(configfields.Args[i])
			}
		} else {
			m1[i] = Prompt(configfields.Args[i])
		}
		newsplice = append(newsplice, m1[i]+"::::")
	}

	bar(secustom)
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
		os.Exit(1)
	}
	out := make([]byte, nonceSize)
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce, naclKey)
	err = ioutil.WriteFile(os.Getenv("HOME")+"/."+secustom, out, 0600)
	if err != nil {
		fmt.Println("Error while writing config file: ", err)
		os.Exit(1)
	}
	fmt.Printf("Config file saved at "+os.Getenv("HOME")+"/."+secustom+" \nTotal size is %d bytes.\n",
		len(out))
	os.Exit(0)
}

func Detect(secustom string) bool {
	_, err := ioutil.ReadFile(os.Getenv("HOME") + "/." + secustom)
	if err != nil {
		return false
	}
	return true
}

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
	in, err := ioutil.ReadFile(os.Getenv("HOME") + "/." + secustom)
	if err != nil {
		fmt.Println(err)

	}
	copy(nonce[:], in[:nonceSize])
	configbytes, ok := secretbox.Open(nil, in[nonceSize:], nonce, naclKey)
	if !ok {
		fmt.Println("Could not decrypt the config file. Wrong password?")
		os.Exit(1)
	}
	//configstrings := strings.Split(string(configbytes), "::::")

	//	username = configstrings[0]
	//	password = configstrings[1]

	return string(configbytes), nil

}
func bar(secustom string) {
	versionbar := strings.Repeat("#", 10) + "\t" + secustom + "\t" + strings.Repeat("#", 30)
	print("\033[H\033[2J")
	fmt.Println(versionbar)
}
