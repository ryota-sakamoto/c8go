build:
	go build -o c8go *.go

test:
	make build
	./test.sh
	make clean

clean:
	rm a.s a *.c
