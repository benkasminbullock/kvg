BINARIES=\
bogusgroup \
missing-stroke \
read-write-test \
renumber \
skip \


all: $(BINARIES)

missing-stroke: missing-stroke.go
	go build missing-stroke.go

read-write-test: read-write-test.go
	go build read-write-test.go

renumber: renumber.go
	go build renumber.go

skip: skip.go
	go build skip.go

bogusgroup: bogusgroup.go
	go build bogusgroup.go

skip.json: make-skip-json.pl
	perl make-skip-json.pl

clean:
	rm -f $(BINARIES)
