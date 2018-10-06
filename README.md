# SafePrice
SafePrice is an applet for checking SafeCoin price
## Build
First of all, you need Golang installed in your system. If it isn't, check this guide: https://golang.org/doc/install
##### Dependency
`go get github.com/getlantern/systray`
### Windows
`go build -ldflags "-s -w -H=windowsgui" main.go`
### Linux & Mac
`go build -ldflags "-s -w" main.go`
