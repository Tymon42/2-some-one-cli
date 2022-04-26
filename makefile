RUMDIR = quorum
GOARCH = amd64
QUORUM_BIN_NAME=quorum
APPDIR = 2SOMEone
LIBVLC = libvlc

# cur_mkfile := $(abspath $(lastword $(MAKEFILE_LIST)))
# cur_makefile_path=$(dir $(cur_mkfile))
# $(info ${cur_makefile_path})

2SOMEone:
	mkdir ${APPDIR}

Copylibvlc:
	Xcopy $(LIBVLC)/ ./$(APPDIR)/ /E/H/C/I

quorum.exe:
	cd $(RUMDIR) && make windows -B
	cd ..
	cp $(RUMDIR)/dist/windows_${GOARCH}/${QUORUM_BIN_NAME}.exe ./${APPDIR}/${QUORUM_BIN_NAME}.exe

# run-peer.exe:
# 	cd ./run-peer && go build -o ../${APPDIR}/run-peer.exe
# 	cd ..

# cli.exe:export CGO_LDFLAGS=-L${cur_makefile_path}/${LIBVLC}
# cli.exe:export CGO_CFLAGS=-I${cur_makefile_path}/${LIBVLC}/include
cli.exe:
	go build -o ./${APPDIR}/cli.exe main.go

gui-nocmd.exe:
	go build -ldflags="-H windowsgui" -o ./${APPDIR}/gui.exe main.go

gui.exe:
	go build -o ./${APPDIR}/gui.exe main.go

#go build -ldflags '-s -w -L C:/libvlc -linkmode "external" -extldflags "-static"' -o ./${APPDIR}/gui.exe main.go

build-gui: Copylibvlc gui.exe

build-win: 2SOMEone quorum.exe cli.exe

clean-win:
	cd ${APPDIR} && rm -rf *
	# rm -rf quorum.exe cli.exe

clean-build-win:clean-win build-win