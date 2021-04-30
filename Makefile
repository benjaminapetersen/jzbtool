all: build
.PHONY: all

OUT_DIR=binaries


# Remove all build artifacts.
#
# Example:
#   make clean
clean:
	rm -rf $(OUT_DIR)
.PHONY: clean

# TODO: left this undone.
# need to get bavk to it
build:
	cd ./src 
	# for linux:
	#  env GOOS=linux GOARCH=arm go build
	# for osx, go build is sufficient.
	go build -o ../binaries/osx/jzbtool
.PHONY: build