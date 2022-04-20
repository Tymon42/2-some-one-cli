# Build on windows  
In the makefile dir:
```bash
$ set CGO_LDFLAGS=-L<libvlc absolutely dir>
$ set CGO_CFLAGS=-I<libvlc absolutely dir>\include
$ make build-win
```

Click `cli.exe` to run.  

Defalt rum api url: `https://127.0.0.1:8004`