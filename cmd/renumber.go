package main

import (
	"kvg"
	"os"
)

/* Renumber all the groups and paths in a file consecutively. */

func main() {
	for i, file := range os.Args {
		if i == 0 {
			continue
		}
		kvg.RenumberFile(file)
	}
}
