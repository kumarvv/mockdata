go mod tidy

echo "### mac"
export CGO_ENABLED=0
export GOOS=darwin
export GOARCH=amd64
go build -o dist/mockdata_mac

echo "### linux"
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o dist/mockdata_linux

echo "### windows"
export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64
go build -o dist/mockdata.exe
