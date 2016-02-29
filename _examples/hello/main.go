package main

import (
"github.com/aerth/seconf"
"os"
"fmt"
"strings"
)
func main() {
  if len(os.Args) < 3 {
      fmt.Println("Usage:")
      fmt.Println(os.Args[0]+" configname servicename field1 field2 etc")
      os.Exit(1)
  }

s := os.Args[1]
sn := os.Args[2]
var fields []string
fields = os.Args[3:]
//fmt.Println(fields)
//os.Exit(1)
if !seconf.Detect(s) {

  seconf.Create(s, sn, fields...)
  //fmt.Println("seconf.Create(s, sn, fields...)")

  }else{

    configdecoded, err := seconf.Read(s)
    if err != nil {
      fmt.Println("error:")
      fmt.Println(err)
      os.Exit(1)
    }
		configarray := strings.Split(configdecoded, "::::")
    //fmt.Println(configarray)
    if len(configarray) < 2 {
      fmt.Println("Broken config file. Create a new one.")
      os.Exit(1)
    }
    //u, p, err := seconf.Read(s)

    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    fmt.Println("Welcome to "+sn+", "+configarray[0])
    fmt.Printf("Your %s is %s \n", os.Args[3], configarray[0])
    fmt.Printf("Your %s is %s \n", os.Args[4], configarray[1])
    fmt.Printf("Your %s is %s \n", os.Args[5], configarray[2])
  //  fmt.Printf("Your %s is %s \n", os.Args[4], configarray[3])

    //fmt.Println("Your password is: "+p)


  }
}
