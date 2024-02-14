.PHONY: clean build zip

GOOS=linux
GOARCH=amd64

clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf bootstrap.zip

build: clean
	@echo "Building..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build  -o ./bin/bootstrap  ./cmd/*.go
	cd bin && chmod +x bootstrap && cd ..


zip: build
	@echo "Zipping..."
	cd bin && zip ../bootstrap.zip bootstrap && cd ..
	clear


all: zip
