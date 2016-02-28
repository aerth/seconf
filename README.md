# Work in progress.

This library allows non-plaintext configuration files

copyright (c) 2016 aerth@sdf.org

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
