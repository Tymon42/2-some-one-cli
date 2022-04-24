package util

import (
	// "2-some-one-cli/util"
	"fmt"
	"os"
	"os/exec"

	// "io/ioutil"
	// "strconv"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	root    = filepath.Join(filepath.Dir(b), "../")

// basepath = filepath.Dir(b)
)

const (
	KeystorePassword = "123"
)

func RunPeer() {

	// pidch := make(chan int)
	// stdoutline := make(chan string)
	root := filepath.Join(root, "2SOMEone")
	// root := ".\\"
	peername := "user2"
	confdir := fmt.Sprintf("%s\\%s", root, "config")
	userdatadir := fmt.Sprintf("%s\\%s", root, "data")
	peerkeystoredir := fmt.Sprintf("%s\\%s", root, peername+"keystore")
	peertracer := fmt.Sprintf("%s\\%s", root, peername+"tracer.json")

	peerport := 7004
	peerapiport := 8004
	// peer := "/ip4/47.96.96.172/tcp/10666/p2p/16Uiu2HAmTHX2YpRTZs4BwcXACGTvJoCRpYWa7paereACj9dduemJ"
	// peer := "/ip4/47.96.96.172/tcp/10666/p2p/16Uiu2HAkxRLEaRonNf4eVdNGdwxE5nh4vGNAUUbjaH8K8HR3yBsC"
	peer := "/ip4/94.23.17.189/tcp/10666/p2p/16Uiu2HAmGTcDnhj3KVQUwVx8SGLyKBXQwfAxNayJdEwfsnUYKK4u"
	// util.Fork(pidch, stdoutline, KeystorePassword, ".\\quorum.exe","-peername", peername, "-listen", fmt.Sprintf("%d%d", "/ip4/127.0.0.1/tcp/", peerport), "-apilisten", fmt.Sprintf("%d", peerapiport), "-peer", peer, "-configdir", confdir, "-datadir", userdatadir, "-keystoredir", peerkeystoredir, "-jsontracer", peertracer)
	cmd := exec.Command("quorum.exe", "-peername", peername, "-listen", fmt.Sprintf("%s%d", "/ip4/127.0.0.1/tcp/",peerport), "-apilisten", fmt.Sprintf(":%d", peerapiport), "-peer", peer, "-configdir", confdir, "-datadir", userdatadir, "-keystoredir", peerkeystoredir, "-jsontracer", peertracer)

	fmt.Println(cmd.Args)

	cmd.Env = append(os.Environ(),
		"RUM_KSPASSWD="+KeystorePassword,
	)

	cmd.Run()
	// cmd.Start()

}
