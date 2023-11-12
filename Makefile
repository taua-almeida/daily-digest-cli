# The name of your binary
BINARY_NAME=gh-digest

# Default command
all: windows linux mac

# Commands for each operating system
windows:
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}.exe main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-linux main.go

mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}-mac main.go

# Clean up binaries
clean:
	rm -f bin/${BINARY_NAME}*
