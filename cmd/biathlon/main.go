package main

import (
	"flag"
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

	config := flag.String("config", "configs/config.json", "path to json contains Competition parameters (configuration file)")
	events := flag.String("events", "test/integration/testdata/events", "path to file contains a set of external events of certain format")
	flag.Parse()

	var compete biathlon
	compete.participants = make(map[int]stat)
	data, err := readJSON(*config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wrong config file (json)")
		os.Exit(-1)
	}
	compete.countLaps = data.Laps
	compete.startDelta, _ = parseTime(data.StartDelta)
	err = proc(&compete, *events)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wrong events file")
		os.Exit(-1)
	}
	output(&compete, &data)
}
