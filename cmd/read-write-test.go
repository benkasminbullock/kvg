/*
   Test reading and then writing an XML file.

   The purpose of this is to see whether we can read in one of the
   KanjiVG inputs and write it out in a similar format, so that the
   git diff does not get messed up when automatic editing is performed
   on the characters.

   This was one of the first things I wrote using the kvg library, and
   thus some of the methods used predate better methods I invented
   later after getting experience with the library. Thus this file may
   not show the best possible practices for using the library.
*/

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"kvg"
	"os"
	"path/filepath"
	"strings"
)

// True if c is whitespace.
func white(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

// Compare two xmls and print the first difference to stdout.
func compareXML(file string, xmlout, xmlin []byte) {
	fileprinted := false
	if len(xmlout) != len(xmlin) {
		fmt.Printf("%s:\n", file)
		fileprinted = true
		fmt.Printf("Lengths differ (%d - %d).\n", len(xmlout), len(xmlin))
	}
	n := len(xmlin)
	if len(xmlout) < n {
		n = len(xmlout)
	}
	line := 1
	offset := 0
	/* The byte where the line starts. */
	lineStart := 0
	for i, c := range xmlin {
		if i >= n {
			break
		}
		if xmlout[i] != c {
			totalFails++
			if white(xmlout[i]) || white(c) {
				fmt.Printf("%s: Whitespace difference.\n", file)
				whiteFails++
			} else {
				fmt.Printf("%s: Attribute or other difference.\n", file)
			}
			if !fileprinted {
				fmt.Printf("%s:\n", file)
				fileprinted = true
			}
			fmt.Printf("First difference at byte %d, line %d, offset %d\n",
				i, line, offset)
			start := lineStart
			end := lineStart + offset + 40
			if end > n {
				end = n
			}
			fmt.Printf("IN:  *%s*\nOUT: *%s*\n", xmlin[start:end], xmlout[start:end])
			if fix {
				// Write the file back out
				err := ioutil.WriteFile(file, xmlout, 0644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing %s: %s\n", file, err)
					os.Exit(1)
				}
			}
			break
		}
		offset++
		if c == '\n' {
			line++
			offset = 0
			lineStart = i + 1
		}
	}
}

var kanjiRad map[rune]map[string]string

func el(g *kvg.Group) (element string) {
	element = g.Element
	if len(g.Original) > 0 {
		element = g.Original
	}
	return element
}

func checkABoo(file string, kanji rune, what string, gs []*kvg.Group) {
	if len(gs) == 0 {
		return
	}
	gen := kanjiRad[kanji][what]
	if len(gen) == 0 {
		if len(kanjiRad[kanji]) == 0 {
			kanjiRad[kanji] = make(map[string]string, 0)
		}
		kanjiRad[kanji][what] = el(gs[0])
		return
	}
	for _, g := range gs {
		el := el(g)
		if el != gen {
			fmt.Printf("%s: %s radical does not match other variant files '%s' (%X) != '%s' (%X)\n",
				file, what, el, []rune(el)[0], gen, []rune(gen)[0])
		}
	}
}

func checkSame(file string, kanji rune, rad kvg.Radical) {
	checkABoo(file, kanji, "general", rad.General)
	checkABoo(file, kanji, "nelson", rad.Nelson)
	checkABoo(file, kanji, "tradit", rad.Tradit)
	checkABoo(file, kanji, "jis", rad.JIS)
}

func checkRadical(file string, svg *kvg.SVG, base *kvg.Group, kanji rune) {
	if !kvg.ExpectRadical(kanji) {
		return
	}
	var rad kvg.Radical
	base.SearchRadical(&rad)
	if len(rad.General) == 0 && len(rad.Tradit) == 0 &&
		len(rad.Nelson) == 0 && len(rad.JIS) == 0 {
		fmt.Printf("No radical found in %s\n", file)
		totalFails++
	}
	if len(rad.Nelson) > 0 && len(rad.Tradit) == 0 {
		fmt.Printf("Inconsistent radicals Nelson, no Tradit in %s\n", file)
		totalFails++
	}
	checkSame(file, kanji, rad)
}

func readWriteTest(file string) {
	contents, oerr := ioutil.ReadFile(file)
	if oerr != nil {
		fmt.Fprintf(os.Stderr, "Error opening %s: %s\n", file, oerr)
		os.Exit(1)
	}
	svg, err := kvg.ParseKanji(contents)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Too bad %s\n", err)
		os.Exit(1)
	}
	id, kanji, _ := kvg.FileToParts(file)
	_, base := svg.Base()
	baseGroup := svg.BaseGroup()
	baseElement := baseGroup.Element
	rebased := false
	if len(baseElement) > 0 {
		baseKanji := []rune(baseElement)[0]
		if int64(baseKanji) != kanji {
			abbrevname := strings.Replace(file, kvg.KVDir+"/", "", 1)
			fmt.Printf("File name, %c, [%s] disagrees with element %s [%05x]\n",
				rune(kanji), abbrevname, baseElement, int64(baseKanji))
			if fix {
				baseGroup.Element = string([]rune{rune(kanji)})
			}
		}
	}
	if svg.Groups[0].ID != "kvg:StrokePaths_"+id {
		fmt.Printf("StrokePaths id %s != %s\n", svg.Groups[0].ID, id)
		totalFails++
		if fix && !rebased {
			svg.SetBase("kvg:" + id)
			rebased = true
		}
	}
	if len(svg.Groups) > 1 && svg.Groups[1].ID != "kvg:StrokeNumbers_"+id {
		fmt.Printf("StrokeNumbers id %s != %s\n", svg.Groups[1].ID, id)
		totalFails++
		if fix && !rebased {
			svg.SetBase("kvg:" + id)
			rebased = true
		}
	}
	if id != base {
		fmt.Printf("Error: base name '%s' and file ID '%s' differ.\n",
			base, id)
		totalFails++
		if fix && !rebased {
			svg.SetBase("kvg:" + id)
			rebased = true
		}
	}
	checkRadical(file, &svg, baseGroup, rune(kanji))
	if len(baseGroup.Position) != 0 {
		fmt.Printf("%s: base group has silly position %s\n",
			file, baseGroup.Position)
		if fix {
			baseGroup.Position = ""
		}
		totalFails++
	}
	svg.RenumberXML()
	xmlout := svg.MakeXML()
	compareXML(file, xmlout, contents)
}

var fix = false
var totalFails = 0
var whiteFails = 0

func main() {
	kanjiRad = make(map[rune]map[string]string, 0)
	fixFlag := flag.Bool("fix", false, "Fix the errors found")
	flag.Parse()
	fix = *fixFlag
	filepath.WalkDir(kvg.KVDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if kvg.Backup.MatchString(path) {
			return nil
		}
		readWriteTest(path)
		return nil
	})
	fmt.Printf("Total failures %d\n", totalFails)
	fmt.Printf("Whitespace-only inconsistencies %d\n", whiteFails)
}
