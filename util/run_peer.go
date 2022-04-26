package util

import (
	// "2-some-one-cli/util"
	"fmt"
	"os"
	"os/exec"

	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	root    = filepath.Join(filepath.Dir(b), "../")
)

func RunPeer(peername string, peerport, peerapiport int, keystore_pw string, peer string) {

	// pidch := make(chan int)
	// stdoutline := make(chan string)
	root := filepath.Join(root, "2SOMEone")
	// root := ".\\"
	// peername := "user2"
	confdir := fmt.Sprintf("%s\\%s", root, "config")
	userdatadir := fmt.Sprintf("%s\\%s", root, "data")
	peerkeystoredir := fmt.Sprintf("%s\\%s", root, peername+"keystore")
	peertracer := fmt.Sprintf("%s\\%s", root, peername+"tracer.json")

	// peer := "/ip4/94.23.17.189/tcp/10666/p2p/16Uiu2HAmGTcDnhj3KVQUwVx8SGLyKBXQwfAxNayJdEwfsnUYKK4u"
	// util.Fork(pidch, stdoutline, KeystorePassword, ".\\quorum.exe","-peername", peername, "-listen", fmt.Sprintf("%d%d", "/ip4/127.0.0.1/tcp/", peerport), "-apilisten", fmt.Sprintf("%d", peerapiport), "-peer", peer, "-configdir", confdir, "-datadir", userdatadir, "-keystoredir", peerkeystoredir, "-jsontracer", peertracer)
	cmd := exec.Command("quorum.exe", "-peername", peername, "-listen", fmt.Sprintf("%s%d", "/ip4/127.0.0.1/tcp/",peerport), "-apilisten", fmt.Sprintf(":%d", peerapiport), "-peer", peer, "-configdir", confdir, "-datadir", userdatadir, "-keystoredir", peerkeystoredir, "-jsontracer", peertracer)

	fmt.Println(cmd.Args)

	cmd.Env = append(os.Environ(),
		// "RUM_KSPASSWD="+KeystorePassword,
		"RUM_KSPASSWD="+keystore_pw,
	)

	cmd.Run()
	// cmd.Start()

}
