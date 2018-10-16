# SafePrice
SafePrice is an applet for checking SafeCoin price
## Usage
**Run it.**

It actually supports some cli arguments. You can check them with `./safeprice --help`
## Build
First of all, you need Golang installed in your system. If it isn't, check this guide: https://golang.org/doc/install
##### Dependency
`go get github.com/getlantern/systray`
### Windows
`go build -ldflags "-s -w -H=windowsgui" safeprice.go`
### Linux & Mac
`go build -ldflags "-s -w" safeprice.go`
