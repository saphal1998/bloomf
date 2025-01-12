default: test 

clean:
	rm -rf ./bloomf

build: clean
	go build -o ./bloomf main.go

test:
	go test ./...
