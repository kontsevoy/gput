package main

import (
	"fmt"
	_ "io"
	"io/ioutil"
	"os"
	"text/scanner"
)

type IniSection struct {
}

type IniConfig struct {
	m map[string]IniSection
}

func ParseIniFile(fileName string) error {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	// use scanner:
	var s scanner.Scanner
	file, err := os.Open(fileName)
	s.Init(file)

	tok := s.Scan()
	for tok != scanner.EOF {
		fmt.Println(s.TokenText())
		tok = s.Scan()
	}

	return nil
}
