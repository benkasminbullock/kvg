BINARIES=\
missing-stroke \
read-write-test \
renumber \


all: $(BINARIES)

missing-stroke: missing-stroke.go
	go build missing-stroke.go

read-write-test: read-write-test.go
	go build read-write-test.go

renumber: renumber.go
	go build renumber.go

clean:
	rm -f $(BINARIES)
