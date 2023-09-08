build:
	go build -o target/server server.go routes.go

clean:
	rm -f target/server
