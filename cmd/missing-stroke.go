package main

import (
	"fmt"
	"kvg"
	"strings"
)

func xm(file string) {
	_, base := kvg.Grab(file)
	paths := base.GetPaths()
	missing := false
	for _, p := range paths {
		if p.Type == "Missing stroke" {
			missing = true
		}
	}
	if missing {
		file = strings.Replace(file, "/home/ben/software/kanjivg/kanji/", "", -1)
		fmt.Printf("%s\n", file)
	}
}

func main() {
	kvg.ExamineAllFilesSimple(xm)
}
