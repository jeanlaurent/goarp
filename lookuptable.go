package main

import (
	"io/ioutil"
	"strings"
)

type LookupTable map[string]string

func newLookupTable() LookupTable {
	return map[string]string{}
}

func (l *LookupTable) loadFromFile(filename string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		tokens := strings.Split(line, ">")
		if len(tokens) >= 2 {
			(*l)[tokens[0]] = tokens[1]
		}
	}
	return nil
}

func (l *LookupTable) get(key string) string {
	return l.get(key)
}
