build:
	GOOS=windows GOARCH=386 go build -o supersafe.exe main.go
	cp supersafe.exe ~/vm-share
