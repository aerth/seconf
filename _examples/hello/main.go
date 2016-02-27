package main

import (
"github.com/aerth/seconf"
"os"
"fmt"
)
func main() {
  if len(os.Args) < 2 {
      fmt.Println("Usage:")
      fmt.Println(os.Args[0]+" configname")
      os.Exit(1)
  }

s := os.Args[1]

if !seconf.Detect(s) {

  seconf.Create(s)

  }else{

    u, _, err := seconf.Read(s)

    //u, p, err := seconf.Read(s)

    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    fmt.Println("Welcome Back, "+u)

    //fmt.Println("Your password is: "+p)


  }
}
