package dockerfile

import (
	"io"
	"os"
	"sort"

	"github.com/docker/docker/builder/dockerfile/command"
	"github.com/docker/docker/builder/dockerfile/parser"
)

// Represents a single line (layer) in a Dockerfile.
// For example `FROM ubuntu:xenial`
type Command struct {
	Cmd       string   // lowercased command name (ex: `from`)
	SubCmd    string   // for ONBUILD only this holds the sub-command
	Json      bool     // whether the value is written in json form
	Original  string   // The original source line
	StartLine int      // The original source line number
	Flags     []string // Any flags such as `--from=...` for `COPY`.
	Value     []string // The contents of the command (ex: `ubuntu:xenial`)
}

// A failure in opening a file for reading.
type IOError struct {
	Msg string
}

func (e IOError) Error() string {
	return e.Msg
}

// A failure in parsing the file as a dockerfile.
type ParseError struct {
	Msg string
}

func (e ParseError) Error() string {
	return e.Msg
}

// List all legal cmds in a dockerfile
func AllCmds() []string {
	var ret []string
	for k := range command.Commands {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

// Parse a Dockerfile from a reader.  A ParseError may occur.
func ParseReader(file io.Reader) ([]Command, error) {
	directive := parser.Directive{LookingForDirectives: true}
	parser.SetEscapeToken(parser.DefaultEscapeToken, &directive)
	ast, err := parser.Parse(file, &directive)
	if err != nil {
		return nil, ParseError{err.Error()}
	}

	var ret []Command
	for _, child := range ast.Children {
		cmd := Command{
			Cmd:       child.Value,
			Original:  child.Original,
			StartLine: child.StartLine,
			Flags:     child.Flags,
		}

		// Only happens for ONBUILD
		if child.Next != nil && len(child.Next.Children) > 0 {
			cmd.SubCmd = child.Next.Children[0].Value
			child = child.Next.Children[0]
		}

		cmd.Json = child.Attributes["json"]
		for n := child.Next; n != nil; n = n.Next {
			cmd.Value = append(cmd.Value, n.Value)
		}

		ret = append(ret, cmd)
	}
	return ret, nil
}

// Parse a Dockerfile from a filename.  An IOError or ParseError may occur.
func ParseFile(filename string) ([]Command, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, IOError{err.Error()}
	}
	defer file.Close()

	return ParseReader(file)
}
