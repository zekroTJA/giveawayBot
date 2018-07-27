package main

import (
	"fmt"
	"strings"
)

//////////////// PUBLIC STUFF ////////////////

// LogInfo crates info log in console
func LogInfo(stuff ...string) {
	for _, e := range(joinAndSplit(stuff)) {
		fmt.Println("[ INFO ] " + e)
	}
}

// LogWarn crates warn log in console
func LogWarn(stuff ...string) {
	for _, e := range(joinAndSplit(stuff)) {
		fmt.Println("[ WARN ] " + e)
	}
}

// LogError crates error log in console
func LogError(stuff ...string) {
	for _, e := range(joinAndSplit(stuff)) {
		fmt.Println("[ ERR! ] " + e)
	}
}

// CheckError checks if error is nil, if not,
// panic if soft is set to false, else log
// error in console
func CheckError(err error, soft bool) {
	if err != nil {
		if !soft {
			panic(err)
		}
		LogError(err.Error())
	}
}


//////////////// RIVATE STUFF ////////////////

func joinAndSplit(stuff []string) []string {
	s := strings.Join(stuff, " ")
	return strings.Split(s, "\n")
}