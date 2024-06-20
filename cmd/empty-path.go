// Search all the files for paths which are empty.

package main

import (
	"fmt"
	"kvg"
)

func emptyPath(file string) {
	svg, base := kvg.Grab(file)
	paths := base.GetPaths()
	nums := &svg.Groups[1]
	nc := len(nums.Children)
	if nc == len(paths) {
		// This file is OK, the number of paths is the same as the number
		// of stroke labels.
		return
	}
	if nc < len(paths) {
		fmt.Printf("%s: missing %d numbers.\n", kvg.TFile(file), len(paths)-nc)
	}
	if nc > len(paths) {
		// This does not happen for any file.
		fmt.Printf("%s: too many stroke numbers %d > %d.\n",
			kvg.TFile(file), nc, len(paths))
		return
	}
	emptyPaths := make([]bool, len(paths))
	found := false
	for i, p := range paths {
		if i < nc {
			continue
		}
		if len(p.D) == 0 {
			emptyPaths[i] = true
			found = true
		}
	}
	if !found {
		return
	}
	newchild := make([]kvg.Child, 0)
	for i := range nums.Children {
		if !emptyPaths[i] {
			newchild = append(newchild, nums.Children[i])
		}
	}
	nums.Children = newchild
	fmt.Printf("%s\n", kvg.TFile(file))
}

func main() {
	kvg.ExamineAllFilesSimple(emptyPath)
}
