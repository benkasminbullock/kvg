package main

import (
	"fmt"
	"kvg"
)

func main() {
	kvg.ExamineAllFilesSimple(bogusGroup)
}
func bogusGroup(file string) {
	_, base := kvg.Grab(file)
	groups := base.GetGroups()
	for _, group := range groups {
		if len(group.Element) > 0 {
			continue
		}
		paths := group.GetPaths()
		if len(paths) > 1 {
			continue
		}
		pos := group.Position
		if len(pos) > 0 {
			continue
		}
		rad := group.Radical
		if len(rad) > 0 {
			continue
		}
		phon := group.Phon
		if len(phon) > 0 {
			continue
		}
		fmt.Printf("%s: %s has one child and no element or position\n",
			file, group.ID)
	}
}
