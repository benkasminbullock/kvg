all: read-write-test renumber

renumber: renumber.go
	go build renumber.go

read-write-test: read-write-test.go
	go build read-write-test.go

clean:
	rm -f renumber read-write-test
