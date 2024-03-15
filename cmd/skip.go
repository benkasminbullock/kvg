/* Compare SKIP codes from the data to ones calculated from the
   KanjiVG data. */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kvg"
	"os"
	"regexp"
	"strconv"
	"unicode"
)

var skipdic map[string]string

type skipcode struct {
	shape int // first key
	a, b  int // second and third keys
}

var skipre = regexp.MustCompile("([1-4])-([0-9]+)-([0-9]+)")

var total = 0
var okskip = 0

// It was not possible to make a guess, based on anything.
var guessfail = 0

// Our guess based on the position value was wrong.
var guesswrong = 0

// There is no group so we cannot get the position from it.
var nogroup = 0

// There is a group but it has no position.
var nopos = 0

// There is only one child of the base element.
var oneChild = 0

var oka = 0
var okb = 0

const noGroup = -2
const noPosition = -4
const singleChild = -3

func skipToNums(skip string) (sc skipcode) {
	matches := skipre.FindStringSubmatch(skip)
	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "Failed to match %s\n", skip)
		os.Exit(1)
	}
	sc.shape, _ = strconv.Atoi(matches[1])
	sc.a, _ = strconv.Atoi(matches[2])
	sc.b, _ = strconv.Atoi(matches[3])
	return sc
}

var PrintSingles = false
var PrintNoPos = false
var PrintAMistake = false // print a count mistakes
var PrintWrong = true
var PrintFails = false
var PrintCounts = false

const unknown = -1

func guessShape(base *kvg.Group) (first, second, third int) {
	kanji := base.Element
	if len(base.Children) == 1 {
		if PrintSingles {
			fmt.Printf("%s: Single child\n", kanji)
		}
		oneChild++
		return 4, unknown, unknown
		//		return singleChild, unknown, unknown
	}
	if !base.Children[0].IsGroup {
		if false {
			fmt.Printf("%s: First child is not a group.\n", kanji)
		}
		nogroup++
		return 4, unknown, unknown
		//		return noGroup, unknown, unknown
	}
	child0 := base.Children[0].Group
	pos := child0.Position
	element := child0.Element
	child0paths := child0.GetPaths()
	nchild0 := len(child0paths)
	basepaths := base.GetPaths()
	nbase := len(basepaths)
	nremaining := nbase - nchild0
	switch element {
	case "匚", "囗":
		nchild0++
		nremaining--
	}
	if false {
		// Sometimes skip is 3 and we are 4 and sometimes skip is 4
		// and we are 4. The latter case outnumbers the former.
		if base.Children[1].IsGroup {
			el1 := base.Children[1].Group.Element
			if el1 == "辶" {
				nremaining--
				fmt.Printf("%s %d\n", kanji, nremaining)
			}
		}
	}
	switch pos {
	case "left":
		return 1, nchild0, nremaining
	case "tare", "nyo", "kamae", "⿵A":
		if pos == "kamae" && element == "行" {
			return 1, nchild0, nremaining
		}
		if pos == "tare" && (element == "户" || element == "戸") {
			return 2, nchild0, nremaining
		}
		return 3, nchild0, nremaining
	case "nyoc", "tarec", "⿶2":
		return 3, nremaining, nchild0
	case "top":
		return 2, nchild0, nremaining
	}
	if len(pos) > 0 {
		/* This does not happen. */
		fmt.Printf("Failed to guess for %s\n", pos)
	} else {
		/* These currently lack a position field for at least some
		   cases in KanjiVG. */
		if element == "尺" || element == "几" || element == "广" ||
			element == "弋" || element == "戈" ||
			element == "耂" {
			return 3, nchild0, nremaining
		}
		// Apel's unusual division into top and bottom of 衣.
		if element == "衣" {
			nremaining = nbase - 2
			return 2, 2, nremaining
		}
		if element == "弍" {
			nremaining = nbase - 3
			return 3, 3, nremaining
		}

		if element == "一" || element == "二" {
			return 4, nbase, 1
		}
		if PrintNoPos {
			fmt.Printf("%s Group on child with element %s, but no position\n",
				kanji, child0.Element)
		}
		nopos++
		return 4, nbase, 4 //-1, unknown, unknown
		//		return noPosition, nchild0, nremaining
	}
	return 4, nbase, 4 //-1, unknown, unknown
}

var disagree = 0
var unusual = 0
var mayBeUsual = 0
var mismatch = 0

func makeSkip(file string) {
	_, kanji, _ := kvg.FileToParts(file)
	k := rune(kanji)
	ks := fmt.Sprintf("%c", k)
	skip, ok := skipdic[ks]
	if !ok {
		return
	}
	if !unicode.In(k, unicode.Han) {
		return
	}
	sc := skipToNums(skip)
	_, base := kvg.Grab(file)
	bshape, a, b := guessShape(base)
	isUnusual := false
	if bshape != sc.shape {
		if sc.shape == 4 {
			unusual++
			isUnusual = true
		} else {
			mayBeUsual++
		}
	}
	if bshape == noGroup {
		if false {
			fmt.Printf("No group %s, skip is %s\n", ks, skip)
		}
	}
	if bshape == noPosition {
		if PrintNoPos {
			if isUnusual {
				fmt.Printf("%s: %s %d %d, unusual skip is %s\n",
					kvg.TFile(file), ks, a, b, skip)
			} else {
				fmt.Printf("%s: no position %s %d %d, skip is %s\n",
					kvg.TFile(file), ks, a, b, skip)
			}
		}
	}
	if bshape == sc.shape {
		okskip++
	} else {
		if bshape > 0 {
			if PrintWrong {
				mismatch++
				fmt.Printf("Mismatch %d: %s SKIP %s != %d-%d-%d\n",
					mismatch, ks, skip, bshape, a, b)
			}
			guesswrong++
		} else {
			if PrintFails {
				fmt.Printf("FAIL for %s SKIP %s != %d-%d-%d\n",
					ks, skip, bshape, a, b)
			}
			switch bshape {
			case singleChild:
				if PrintSingles {
					fmt.Printf("%s: skip = %s\n", kvg.TFile(file), skip)
				}
			}
			guessfail++
		}
	}
	if a != unknown && b != unknown {
		strCount := a + b
		skipCount := sc.a + sc.b
		if bshape == 4 {
			strCount = a
		}
		if sc.shape == 4 {
			skipCount = sc.a
		}
		if strCount != skipCount {
			if PrintCounts {
				fmt.Printf("%c: %s: %s ", kanji, kvg.TFile(file), skip)
				fmt.Printf("Stroke count disagreement: %d != %d\n",
					strCount, skipCount)
			}
			disagree++
		}
	}
	if a == sc.a {
		oka++
	} else if a != unknown {
		if PrintAMistake {
			fmt.Printf("%c: %s: %d != %d (skip)\n",
				kanji, kvg.TFile(file), a, sc.a)
		}
	}
	if b == sc.b {
		okb++
	}
	total++
}

func main() {
	skipdata, err := ioutil.ReadFile("skip.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal(skipdata, &skipdic)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	kvg.ExamineAllFilesSimple(makeSkip)
	fmt.Printf("ok %d guess wrong %d total %d  [should = %d]\n",
		okskip, guesswrong, total, okskip+guesswrong)
	fmt.Printf("No group = %d no position = %d, one child = %d [total = %d]\n",
		nogroup, nopos, oneChild, nogroup+nopos+oneChild)
	fmt.Printf("OK a %d b %d [stroke counts disagree %d]\n",
		oka, okb, disagree)
	fmt.Printf("Of total failed & wrong, skip probably unusual %d, may be usual %d\n",
		unusual, mayBeUsual)
}
