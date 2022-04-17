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
	nodename := "user"
	confdir := fmt.Sprintf("%s/%s",rumdata, "config")
	userdatadir := fmt.Sprintf("%s/%s", rumdata, "data")
	peerkeystoredir := fmt.Sprintf("%s/%s", rumdata, nodename+"keystore")
	peertracer := fmt.Sprintf("%s/%s", rumdata, nodename+"tracer.json")
	peerport := 7004
	peerapiport := 8004
	// peer := "/ip4/47.96.96.172/tcp/10666/p2p/16Uiu2HAmTHX2YpRTZs4BwcXACGTvJoCRpYWa7paereACj9dduemJ"
	peer := "/ip4/47.96.96.172/tcp/10666/p2p/16Uiu2HAkxRLEaRonNf4eVdNGdwxE5nh4vGNAUUbjaH8K8HR3yBsC"
	util.Fork(pidch, out_str, KeystorePassword, ".\\quorum.exe", "-peername", nodename, "-listen", fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", peerport), "-apilisten", fmt.Sprintf(":%d", peerapiport), "-peer", peer, "-configdir", confdir, "-keystoredir", peerkeystoredir, "-datadir", userdatadir, "-jsontracer", peertracer, "-debug=true")

	
}