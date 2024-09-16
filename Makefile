build:
	go build -o ./build/arch-mirror-gen
clean:
	rm -rf ./build
install:
	go build -o /usr/local/bin/arch-mirror-gen
uninstall:
	rm -rf /usr/local/bin/arch-mirror-gen

