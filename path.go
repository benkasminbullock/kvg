package kvg

import (
	"strconv"
	"strings"
)

// structure containing information about path commands as by specification
var allCommands commands

type commands struct {
	all        []string
	parameters map[string]int
	start      string
	end        string
}

func getCommands() commands {
	var parameters = map[string]int{
		"m": 2, "z": 0, "l": 2, "h": 1, "v": 1,
		"c": 6, "s": 4, "q": 4, "t": 2, "a": 7,
	}
	var all []string
	for k := range parameters {
		all = append(all, k)
	}
	return commands{all, parameters, "m", "z"}
}

func (c commands) isCommand(token string) bool {
	for _, command := range c.all {
		if strings.ToLower(token) == command {
			return true
		}
	}
	return false
}

// PathParserError contains errors which have occured when parsing 'd'
// attribute of a path element.
type PathParserError struct {
	msg string
}

func (err PathParserError) Error() string {
	return err.msg
}

// token can contain an operator or an operand as string.
type token struct {
	value    string
	operator bool
}

// Command is a representation of an SVG path command and its parameters.
type Command struct {
	Symbol string
	Params []float64
}

// IsAbsolute returns true if the SVG path command is absolute.
func (c Command) IsAbsolute() bool {
	return c.Symbol == strings.ToUpper(c.Symbol)
}

// Compare compares two commands.
func (c Command) Compare(o Command) bool {
	if c.Symbol != o.Symbol {
		return false
	}
	for i, param := range c.Params {
		if param != o.Params[i] {
			return false
		}
	}
	return true
}

// Subpath is a collection of Commands, beginning with moveto command and
// usually ending with closepath command.
type Subpath struct {
	Commands []Command
}

// Compare compares two subpaths.
func (s Subpath) Compare(o Subpath) bool {
	if len(s.Commands) != len(o.Commands) {
		return false
	}
	for i, command := range s.Commands {
		if !command.Compare(o.Commands[i]) {
			return false
		}
	}
	return true
}

// SVGPath is a collection of all the subpaths in 'd' attribute.
type SVGPath struct {
	Subpaths []Subpath
}

// Compare compares two paths.
func (p SVGPath) Compare(o SVGPath) bool {
	if len(p.Subpaths) != len(o.Subpaths) {
		return false
	}
	for i, subpath := range p.Subpaths {
		if !subpath.Compare(o.Subpaths[i]) {
			return false
		}
	}
	return true
}

func reverse(ops []float64) []float64 {
	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}
	return ops
}

func addOperand(tokens []token, operand string) ([]token, string) {
	if operand != "" {
		tokens = append(tokens, token{operand, false})
		operand = ""
	}
	return tokens, operand
}

// tokenize takes value of 'd' attribute and transforms it to series of
// operators and operands - step 1.
func tokenize(raw string) []token {
	var (
		tokens  []token
		operand string
	)
	for _, r := range raw {
		char := string(r)
		switch {
		case allCommands.isCommand(char):
			tokens, operand = addOperand(tokens, operand)
			tokens = append(tokens, token{char, true})
		case char == ".":
			if operand == "" {
				operand = "0"
			}
			if strings.Contains(operand, char) {
				tokens = append(tokens, token{operand, false})
				operand = "0"
			}
			fallthrough
		case char >= "0" && char <= "9" || char == "e":
			operand += char
		case char == "-":
			if strings.HasSuffix(operand, "e") {
				operand += char
			} else {
				tokens, operand = addOperand(tokens, operand)
				operand = char
			}
		default:
			tokens, operand = addOperand(tokens, operand)
		}
	}
	tokens, operand = addOperand(tokens, operand)
	return tokens
}

// toCommands takes a collection of operators and operands and produces
// Command objects - step 2.
func toCommands(tokens []token) (commands []Command, err error) {
	var operands []float64
	for i := len(tokens) - 1; i >= 0; i-- {
		t := tokens[i]
		if t.operator {
			nParam := allCommands.parameters[strings.ToLower(t.value)]
			nOperand := len(operands)
			if nParam == 0 && nOperand == 0 {
				command := Command{Symbol: t.value}
				commands = append([]Command{command}, commands...)
			} else if nParam != 0 && nOperand%nParam == 0 {
				loopCount := nOperand / nParam
				for i := 0; i < loopCount; i++ {
					operator := t.value
					if operator == "m" && i < loopCount-1 {
						operator = "l"
					}
					if operator == "M" && i < loopCount-1 {
						operator = "L"
					}
					command := Command{operator, reverse(operands[:nParam])}
					commands = append([]Command{command}, commands...)
					operands = operands[nParam:]
				}
			} else {
				err := PathParserError{"Incorrect number of parameters for " + t.value}
				return commands, err
			}
		} else {
			number, err := strconv.ParseFloat(t.value, 64)
			if err != nil {
				return commands, err
			}
			operands = append(operands, number)
		}
	}
	return commands, nil
}

// createSubpaths takes a collection of Command objects and determines
// all subpaths within the collection - step 3.
func createSubpaths(commands []Command) (path SVGPath) {
	var subpath []Command
	for i, command := range commands {
		switch strings.ToLower(command.Symbol) {
		case allCommands.start:
			if len(subpath) > 0 {
				path.Subpaths = append(path.Subpaths, Subpath{subpath})
			}
			subpath = []Command{command}
		case allCommands.end:
			subpath = append(subpath, command)
			path.Subpaths = append(path.Subpaths, Subpath{subpath})
			subpath = []Command{}
		default:
			subpath = append(subpath, command)
			if len(commands) == i+1 {
				path.Subpaths = append(path.Subpaths, Subpath{subpath})
			}
		}
	}
	return path
}

// PathParser takes value of a 'd' attribute and transforms it to collection of
// subpaths and commands.
func PathParser(raw string) (path SVGPath, err error) {
	allCommands = getCommands()
	commands, err := toCommands(tokenize(raw))
	if err != nil {
		return path, err
	}
	return createSubpaths(commands), nil
}
