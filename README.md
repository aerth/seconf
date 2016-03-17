# seconf

### This library creates, detects, and reads non-plaintext configuration files.

[![GoDoc](https://godoc.org/github.com/aerth/seconf?status.svg)](https://godoc.org/github.com/aerth/seconf)

Currently, seconf saves the configuration file as a `::::` separated list. For now, it works. It would probably better to use JSON or something.
I created this for [go-quitter](https://github.com/aerth/go-quitter), so that the username and password (and node) would not be stored in plaintext. If your app can use it, go ahead! Things may change/break. [Cosgo](https://github.com/aerth/cosgo) also uses it.

copyright (c) 2016 aerth@sdf.org


### [(Example)](https://github.com/aerth/seconf/blob/master/_examples/hello/main.go)


