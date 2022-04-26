package main

import (
	"2-some-one-cli/term"
	"2-some-one-cli/util"
	"log"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	PeerConfig peerConfig
}

type peerConfig struct {
	PeerName         string `toml:"peername"`
	PeerPort         int    `toml:"peerport"`
	PeerApiPort      int    `toml:"peerapiport"`
	KeystorePassword string `toml:"keystore_pw"`
	peer             string `toml:"peer"`
}

func main() {
	var config *tomlConfig = &tomlConfig{}
	filePath := "../config.toml"
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		log.Fatalln(err)
	}
	go util.RunPeer(config.PeerConfig.PeerName, config.PeerConfig.PeerPort, config.PeerConfig.PeerApiPort, config.PeerConfig.KeystorePassword, config.PeerConfig.peer)
	term.Home()
}
