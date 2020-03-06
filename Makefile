linux: linux-386 linux-amd64
linux-386:
	GOOS=linux GOARCH=386 go build -a -o release/vault-migrator-linux-386
linux-amd64:
	GOOS=linux GOARCH=amd64 go build -a -o release/vault-migrator-linux-amd64

darwin: darwin-386 darwin-amd64
darwin-386:
	GOOS=darwin GOARCH=386 go build -a -o release/vault-migrator-darwin-386
darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -a -o release/vault-migrator-darwin-amd64

windows: windows-386 windows-amd64
windows-386:
	GOOS=windows GOARCH=386 go build -a -o release/vault-migrator-windows-386.exe
windows-amd64:
	GOOS=windows GOARCH=amd64 go build -a -o release/vault-migrator-windows-amd64.exe

build: linux darwin windows
