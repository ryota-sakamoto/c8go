build:
	go build -o bin/c8go cmd/c8go/main.go

test:
	make build
	./test.sh
	make clean

clean:
	rm -rf bin tmp
