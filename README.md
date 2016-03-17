# seconf

### This library creates, detects, and reads non-plaintext configuration files.

[![GoDoc](https://godoc.org/github.com/aerth/seconf?status.svg)](https://godoc.org/github.com/aerth/seconf)

Currently, seconf saves the configuration file as a "::::" separated list. It would probably better to use JSON or something.
I created this for go-quitter, so that the username and password would not be stored in plaintext. If your app can use it, go ahead! Things may change/break. [Cosgo](https://github.com/aerth/cosgo) also uses it.

copyright (c) 2016 aerth@sdf.org


### Example

```

cd _examples/hello/

GOOS=windows GOARCH=amd64 go build -o windows.exe
GOOS=darwin GOARCH=amd64 go build -o darwin-amd64
GOOS=linux GOARCH=amd64 go build -o linux-amd64
GOOS=linux GOARCH=arm64 go build -o linux-arm64
GOOS=freebsd GOARCH=amd64 go build -o freebsd-amd64
GOOS=netbsd GOARCH=amd64 go build -o netbsd-amd64
GOOS=openbsd GOARCH=amd64 go build -o openbsd-amd64

```
