package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aerth/seconf"
)

func main() {

	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage:")
		fmt.Println(os.Args[0] + " [keyname]")
		fmt.Println("Example:")
		fmt.Println(os.Args[0] + " \"AWS key\"")
		os.Exit(1)
	}

	s := os.Args[1]
	sn := os.Args[1]
	field := os.Args[1]

	if !seconf.Detect(s) {
		seconf.Create(s, sn, field)
	} else {
		configdecoded, err := seconf.Read(s)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		configarray := strings.Split(configdecoded, "::::")
		if len(configarray) < 1 {
			fmt.Println("Broken config file. Create a new one.")
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Welcome to seconf, " + configarray[0])
		fmt.Printf("Your %s is %s \n", os.Args[1], configarray[0])

	}
}
