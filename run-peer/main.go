package main

import (
	"2-some-one-cli/util"
	"fmt"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	rumdata   = filepath.Join(filepath.Dir(b), "../")
	// basepath = filepath.Dir(b)
)

const (
	KeystorePassword = "123"
)

func main() {

pidch := make(chan int)
	out_str := make(chan string)
	rumdata := filepath.Join(rumdata,"2SOMEone", "rumdata")
	peername := "user"
	userconfdir := fmt.Sprintf("%s/%s",rumdata, peername+"config")
	userdatadir := fmt.Sprintf("%s/%s", rumdata, peername+"data")
	peerkeystoredir := fmt.Sprintf("%s/%s", rumdata, peername+"keystore")
	peerport := 7004
	peerapiport := 8004
	peer := "/ip4/127.0.0.1/tcp/10666/p2p/16Uiu2HAm7YBSxj7Gsw5xz6Ndp5yq28Ku7h6bzJJ5ignk1s6c7o9e"

	util.Fork(pidch, out_str, KeystorePassword, ".\\quorum.exe", "-peername", peername, "-listen", fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", peerport), "-apilisten", fmt.Sprintf(":%d", peerapiport), "-peer", peer, "-configdir", userconfdir, "-keystoredir", peerkeystoredir, "-datadir", userdatadir)

	
}