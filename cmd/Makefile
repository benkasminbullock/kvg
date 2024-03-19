BINARIES=\
bogusgroup \
missing-stroke \
read-write-test \
renumber \
skip \
typeshift \


all: $(BINARIES)

bogusgroup: $@.go
	go build $@.go

missing-stroke: $@.go
	go build $@.go

read-write-test: $@.go
	go build $@.go

renumber: $@.go
	go build $@.go

skip: $@.go
	go build $@.go

typeshift: $@.go
	go build $@.go

skip.json: make-skip-json.pl
	perl make-skip-json.pl

clean:
	rm -f $(BINARIES)
