package main

import (
	"fmt"
	"strings"
)

//////////////// PUBLIC STUFF ////////////////

func LogInfo(stuff ...string) {
	for _, e := range(joinAndSplit(stuff)) {
		fmt.Println("[ INFO ] " + e)
	}
}

func LogWarn(stuff ...string) {
	for _, e := range(joinAndSplit(stuff)) {
		fmt.Println("[ WARN ] " + e)
	}
}

func LogError(stuff ...string) {
	for _, e := range(joinAndSplit(stuff)) {
		fmt.Println("[ ERR! ] " + e)
	}
}


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