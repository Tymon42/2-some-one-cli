package main

import (
	"2-some-one-cli/term"
	"2-some-one-cli/util"
)


func main() {
	go util.RunPeer()
	term.Home()
}
