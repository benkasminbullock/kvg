package kvg

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func die(err error, format string, a ...any) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, ": %s\n", err)
	os.Exit(1)
}

func bin() string {
	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(filename))
	die(err, "Error getting running directory")
	return dir
}

func read(file string) string {
	b, err := os.ReadFile(file)
	die(err, "Error reading %s", file)
	return string(b)
}

func TestWriteKanjiFile(t *testing.T) {
	bin := bin()
	tdir := bin + "/t/"
	infile := tdir + "08475.svg"
	svg, err := ReadKanjiFile(infile)
	if err != nil {
		t.Errorf("Error reading %s: %s", infile, err)
		return
	}
	ofile := bin + "/testoutput.svg"
	svg.WriteKanjiFile(ofile)
	a := read(infile)
	b := read(ofile)
	if a != b {
		t.Errorf("Contents of %s and %s are different.\n", infile, ofile)
	}
	err = os.Remove(ofile)
	die(err, "Error removing %s", ofile)
}
