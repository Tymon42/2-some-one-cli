# Build on windows  
**Before building project, set the CGO_LDFLAGS and the CGO_CFLAGS environment variables, in order to make the Go build tools aware of the location of the VLC SDK files.**

Get your `libvlc` path,  looks like: `C:\Users\why00\Documents\tools\libvlc`  
Then,  
```
$ set CGO_LDFLAGS=-LC:\Users\why00\Documents\tools\libvlc
$ set CGO_CFLAGS=-IC:\Users\why00\Documents\tools\libvlc\include
$ make build-win
```

Click `run-peer.exe` to run peer.  
Click `cli.exe` to run.  

Defalt rum api url: `https://127.0.0.1:8004`