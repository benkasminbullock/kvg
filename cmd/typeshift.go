package main

import (
	"flag"
	"fmt"
	"kvg"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var verbose = true

func main() {
	fileFlag := flag.String("file", "", "File to read")
	writeFlag := flag.Bool("write", false, "Perform a write operation")
	shiftFlag := flag.String("shift", "", "Shifts to perform")
	swapFlag := flag.String("swap", "", "A swap to perform")
	verboseFlag := flag.Bool("verbose", false, "Switch on debugging")
	flag.Parse()
	verbose = *verboseFlag
	if len(*fileFlag) == 0 {
		fmt.Printf("Specify the file with --file <file>\n")
		return
	}
	file := kvg.KVDir + "/" + *fileFlag
	svg := kvg.ReadKanjiFileOrDie(file)
	paths := svg.GetPaths()
	n := len(paths)
	save := make([]string, n)
	for i, p := range paths {
		save[i] = p.Type
	}
	shift := *shiftFlag
	swap := *swapFlag
	var shifts []int
	if len(shift) != 0 {
		if len(swap) != 0 {
			fmt.Fprintf(os.Stderr, "Choose only one of swap or shift.\n")
			return
		}
		shifts = parseShifts(shift, n)
	} else {
		if len(swap) != 0 {
			shifts = parseSwap(swap, n)
		} else {
			fmt.Printf("No shift or swap supplied. Current values are\n")
			for i := range paths {
				fmt.Printf("%d: %s\n", i+1, save[i])
			}
			return
		}
	}
	if !*writeFlag {
		fmt.Printf("If the following alterations look OK, use --write to write this change.\n")
		for i := range paths {
			if false {
				fmt.Printf("%d %d <%s>\n", i, shifts[i], save[shifts[i]])
			}
			if shifts[i] != i {
				fmt.Printf("%d: %s -> %s\n", i+1, save[i], save[shifts[i]])
			}
		}
	}
	if *writeFlag {
		for i := range paths {
			if shifts[i] != i {
				paths[i].Type = save[shifts[i]]
			}
		}
		svg.WriteKanjiFile(file)
	}
}

var digitRange = "([0-9]+)(?:-([0-9]+))?"
var shiftCommand = regexp.MustCompile(digitRange + "=" + digitRange)
var swapCommand = regexp.MustCompile("([0-9]+)=([0-9]+)")

type shiftRange struct {
	begin, end int
}

type shiftInput struct {
	before, after shiftRange
}

func digitError(d string, err error) {
	fmt.Fprintf(os.Stderr, "Error parsing digits %s: %s\n", err)
	os.Exit(1)
}

func parseDigit(d string) int {
	i, err := strconv.Atoi(d)
	if err != nil {
		digitError(d, err)
	}
	return i
}

func parseSwap(swap string, n int) (shifts []int) {
	d := swapCommand.FindStringSubmatch(swap)
	a := parseDigit(d[1])
	b := parseDigit(d[2])
	shifts = blankShifts(n)
	shifts[a-1] = b - 1
	shifts[b-1] = a - 1
	return shifts
}

func parseRange(d []string) (s shiftRange) {
	var err error
	s.begin, err = strconv.Atoi(d[0])
	if err != nil {
		digitError(d[0], err)
	}
	if len(d[1]) > 0 {
		s.end, err = strconv.Atoi(d[1])
		if err != nil {
			digitError(d[1], err)
		}
	}
	return s
}

func rangeSize(s shiftRange) int {
	if s.end == 0 {
		return 1
	}
	return s.end - s.begin + 1
}

func blankShifts(n int) (shifts []int) {
	shifts = make([]int, n)
	for i := 0; i < n; i++ {
		shifts[i] = i
	}
	return shifts
}

// Parse the command into a set of instructions
func parseShifts(shift string, n int) (shifts []int) {
	shifts = blankShifts(n)
	commands := strings.Split(shift, ",")
	for _, c := range commands {
		var s shiftInput
		d := shiftCommand.FindStringSubmatch(c)
		s.before = parseRange(d[1:3])
		s.after = parseRange(d[3:5])
		if verbose {
			fmt.Printf("%d-%d -> %d-%d\n", s.before.begin, s.before.end, s.after.begin, s.after.end)
		}
		bs := rangeSize(s.before)
		as := rangeSize(s.after)
		if bs != as {
			fmt.Fprintf(os.Stderr, "Sizes of ranges in %s differ, %d != %d\n",
				d[0], bs, as)
			os.Exit(1)
		}
		if bs == 1 {
			shifts[s.before.begin-1] = s.after.begin - 1
			continue
		}
		for i := 0; i < bs; i++ {
			shifts[s.before.begin+i-1] = s.after.begin + i - 1
		}
	}
	if verbose {
		for i := range shifts {
			fmt.Printf("%d -> %d\n", i+1, shifts[i]+1)
		}
	}
	if !checkShifts(shifts, n) {
		os.Exit(1)
	}
	return shifts
}

// Check that the shifts make sense (don't put two things into the
// same place). In mathematical terms, check that "shifts" is a valid
// permutation.
func checkShifts(shifts []int, n int) bool {
	exists := make([]bool, n)
	// Check that nothing is mapped to the same value twice.
	for i, s := range shifts {
		if exists[s] {
			fmt.Printf("Duplicate entry at %d (%d)\n", i+1, s+1)
			return false
		}
		exists[s] = true
	}
	// Check that everything has one value.
	for i, e := range exists {
		if !e {
			fmt.Printf("There is no shift-to value for %d\n", i+1)
			return false
		}
	}
	return true
}
