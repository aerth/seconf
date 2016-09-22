package main

import (
	"fmt"
	"os"

	"github.com/aerth/seconf"
)

func main() {
	filename := "se.conf"
	header := "Seconf"
	// []string of fields, they are both the prompt and key for looking the fields back up.
	fields := map[string]string{
		"name":           "What is your name?",
		"favorite-color": "What is your favorite color?",
		"lol":            "You lol tho?",
		"password":       "What is your password? Will not echo",
	}

	if !seconf.Detect(filename) {
		seconf.LockJSON(filename, header, fields) // Ask user for values
		return
	}

	c, err := seconf.ReadJSON(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%s's favorite color is %q!\n", c.Fields["name"], c.Fields["favorite-color"])
}
