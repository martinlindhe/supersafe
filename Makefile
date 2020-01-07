get-compiler:
	go get golang.org/x/build/version/go1.10.3
	go1.10.3 download

build:
	GOOS=windows GOARCH=386 go1.10.3 build -o binaries/supersafe-win32.exe main.go
	upx -9 binaries/supersafe-win32.exe
