// Parser for ini and conf files.
// Usage
//
package main

import (
	"os"
	. "strings"
	. "text/scanner"
)

const (
	TrimChars   = "\"'"
	CommentChar = ";"
)

// IniConfig type stores all values found in a ini-file
type IniConfig struct {
	m map[string]map[string]string
}

func (conf *IniConfig) Get(section string, name string) string {
	return conf.GetSection(section)[normalize(name)]
}

func (conf *IniConfig) GetSection(section string) map[string]string {
	return conf.m[normalize(section)]
}

// ParseIniFile reads the supplied ini-file and returns a IniConf structure
// Later you can use IniConf.Get("section", "name") to get config values
func ParseIniFile(fileName string) (conf IniConfig, err error) {
	var currentSection, currentName string
	conf.m = make(map[string]map[string]string)

	err = processIniFile(fileName,
		// adds a new section to the conf
		func(section string) {
			currentSection = normalize(section)
		},
		func(name string) {
			currentName = normalize(name)
		},
		// adds a new key/value pair to the current section in conf
		func(value string) {
			if value == "" {
				return
			}
			_, haveSection := conf.m[currentSection]
			if !haveSection {
				conf.m[currentSection] = make(map[string]string)
			}
			conf.m[currentSection][currentName] = Trim(value, TrimChars)
		})
	return
}

// normalize() is called on all section names and argument names, making
// them case-insensitive and space-ignoring
func normalize(key string) string {
	return ToLower(Replace(key, " ", "", -1))
}

// processIniFile() actually scans the file, finding config sections
// and name/value pairs, calling provided callbacks for them
func processIniFile(fileName string,
	addSection func(string),
	addName func(string),
	addValue func(string)) (err error) {
	// possible parser states:
	const (
		StateSection = iota
		StateName
		StateValue
		StateComment
	)

	state := StateName // initially look for opening section
	buffer := ""       // buffer to accumulate tokens
	token := ""        // current token
	line := 0          // keeps track of the last line to detect newlines
	var (
		pos Position
		s   Scanner
	)

	// switches parser state and resets buffer
	flipTo := func(newState int) {
		state = newState
		buffer = ""
	}

	// processes one token when parser is in "parsing section" state
	onSection := func() {
		if token == "]" {
			addSection(buffer)
			flipTo(StateName)
		} else {
			buffer += token
		}
	}

	// processes one token when parser is in "parsing parameter name" state
	onName := func() {
		if token == "[" && buffer == "" {
			flipTo(StateSection)
		} else if token == "=" {
			addName(buffer)
			flipTo(StateValue)
		} else {
			buffer += token
		}
	}

	file, err := os.Open(fileName)
	if err != nil {
		return
	}

	// Scan & tokenize the config file:
	s.Init(file)
	for tok := s.Scan(); tok != EOF; tok = s.Scan() {
		pos = s.Pos()
		token = s.TokenText()
		newline := (pos.Line > line)

		// ignore comments right away:
		if newline && token == CommentChar {
			addValue(buffer)
			flipTo(StateComment)
		} else {
			// wich state is the scanner in?
			switch state {
			case StateSection:
				onSection()
			case StateName:
				onName()
			case StateValue:
				if newline {
					addValue(buffer)
					if token == "[" {
						flipTo(StateSection)
						continue
					} else {
						flipTo(StateName)
					}
				}
				buffer += token
			case StateComment:
				if newline { // comment ended
					flipTo(StateName)
					onName()
				}
			}
		}
		line = pos.Line
	}
	// pick up the accumulated buffer as the last value:
	if state == StateValue {
		addValue(buffer)
	}
	return
}
