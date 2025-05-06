package main

import (
	"fmt"
	"log"
	"os"
)

type interval struct {
	compID int
	end    myTime
}

type biathlon struct {
	countLaps    int
	startDelta   myTime
	queue        []interval
	participants map[int]stat
}

func main() {
	log.SetFlags(0)

	var compete biathlon
	compete.participants = make(map[int]stat)
	data, err := readJSON(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "wrong config file (json)")
		os.Exit(-1)
	}
	compete.countLaps = data.Laps
	compete.startDelta, _ = parseTime(data.StartDelta)
	err = proc(&compete, os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "wrong events file")
		os.Exit(-1)
	}
	output(&compete, &data)
}
