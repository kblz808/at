build:
	go build -ldflags="-s -w" -o ./bin/at .

run: build
	./bin/at

clean:
	rm -f ./at
