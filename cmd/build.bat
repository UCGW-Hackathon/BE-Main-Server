set GO111MODULE=on
go version > version.txt
go mod tidy
go build -o app.exe main.go
echo Done > done.txt
