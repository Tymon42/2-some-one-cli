RUMDIR = quorum
GOARCH = amd64
QUORUM_BIN_NAME=quorum
# LIBVLC = C:\Users\why00\Documents\tools\libvlc
# LIBVLC = C:\Users\why00\Documents\tools\vlc-3.0.16
APPDIR = 2SOMEone

2SOMEone:
	mkdir ${APPDIR}

quorum.exe:
	cd $(RUMDIR) && make windows -B
	cd ..
	cp $(RUMDIR)/dist/windows_${GOARCH}/${QUORUM_BIN_NAME}.exe ./${APPDIR}/${QUORUM_BIN_NAME}.exe

run-peer.exe:
	cd ./run-peer && go build -o ../${APPDIR}/run-peer.exe
	cd ..

cli.exe:
	go build -o ./${APPDIR}/cli.exe main.go

build-win: 2SOMEone quorum.exe run-peer.exe cli.exe


clean-win:
	cd ${APPDIR} && rm -rf *
	# rm -rf quorum.exe cli.exe

clean-build-win:clean-win 2SOMEone quorum.exe run-peer.exe cli.exe