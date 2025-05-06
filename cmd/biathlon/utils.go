package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type jsonRace struct {
	Laps         int    `json:"laps"`
	LapLen       int    `json:"lapLen"`
	PenaltyLen   int    `json:"penaltyLen"`
	FiringRanges int    `json:"firingLines"`
	Start        string `json:"start"`
	StartDelta   string `json:"startDelta"`
}

func readJSON(path string) (data jsonRace, err error) {
	var jsonFile []byte
	jsonFile, err = os.ReadFile(path)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(jsonFile, &data)
	return data, err
}

type statusRacer int

const (
	finished statusRacer = iota
	notStarted
	notFinished
	countOfTargets int = 5
)

type lapStat struct {
	enterTime myTime
	totalTime myTime
}

// statistic about one of competitor
// calculating during events hadling
type stat struct {
	status        statusRacer
	startTime     myTime
	lastStartTime myTime
	totalTime     myTime
	penalty       lapStat
	lapsCounter   int
	laps          []lapStat
	hitStat       struct {
		missed int
		hitted int
		total  int
	}
}

// check that nobody skipped his/her timing
func checkQueue(raceStat *biathlon, current myTime) {
	for len(raceStat.queue) > 0 && raceStat.queue[0].end < current {
		fmt.Printf("[%s] 32 %d", getString(raceStat.queue[0].end), raceStat.queue[0].compID)
		val := raceStat.participants[raceStat.queue[0].compID]
		val.status = notStarted
		raceStat.participants[raceStat.queue[0].compID] = val
		raceStat.queue = raceStat.queue[1:]
	}
}

// handle one event in `args`
func handle(raceStat *biathlon, args ...string) error {
	trimmed := strings.Trim(args[0], "[]")
	time, err := parseTime(trimmed)

	if err != nil {
		return err
	}

	eventID, _ := strconv.Atoi(args[1])
	competitor, _ := strconv.Atoi(args[2])

	checkQueue(raceStat, time)

	switch eventID {
	case 1:
		log.Printf("%s The competitor(%s) registered\n", args[0], args[2])
		raceStat.participants[competitor] = stat{laps: make([]lapStat, raceStat.countLaps)}
	case 2:
		log.Printf("%s The start time for the competitor(%s) was set by a draw to %s\n", args[0], args[2], args[3])
		val := raceStat.participants[competitor]
		scheduledStart, _ := parseTime(args[3])

		raceStat.queue = append(raceStat.queue,
			interval{competitor, scheduledStart + raceStat.startDelta})
		val.laps[val.lapsCounter].enterTime = time

		raceStat.participants[competitor] = val
	case 3:
		log.Printf("%s The competitor(%s) is on the start line\n", args[0], args[2])
	case 4:
		log.Printf("%s The competitor(%s) has started\n", args[0], args[2])
		raceStat.queue = raceStat.queue[1:]
	case 5:
		log.Printf("%s The competitor(%s) is on the firing range(%s)\n", args[0], args[2], args[3])
		val := raceStat.participants[competitor]
		val.hitStat.total += countOfTargets
		val.hitStat.missed += countOfTargets
		raceStat.participants[competitor] = val
	case 6:
		log.Printf("%s The target(%s) has been hit by competitor(%s)\n", args[0], args[3], args[2])
		val := raceStat.participants[competitor]

		val.hitStat.missed--
		val.hitStat.hitted++

		raceStat.participants[competitor] = val
	case 7:
		log.Printf("%s The competitor(%s) left the firing range\n", args[0], args[2])
	case 8:
		log.Printf("%s The competitor(%s) entered the penalty laps\n", args[0], args[2])
		val := raceStat.participants[competitor]

		val.penalty.enterTime = time

		raceStat.participants[competitor] = val
	case 9:
		log.Printf("%s The competitor(%s) left the penalty laps\n", args[0], args[2])
		val := raceStat.participants[competitor]

		val.penalty.totalTime += time - val.penalty.enterTime
		val.penalty.enterTime = 0

		raceStat.participants[competitor] = val
	case 10:
		val := raceStat.participants[competitor]
		log.Printf("%s The competitor(%s) ended the main lap\n", args[0], args[2])

		val.laps[val.lapsCounter].totalTime = time - val.laps[val.lapsCounter].enterTime
		val.lapsCounter++
		if val.lapsCounter == len(val.laps) {
			fmt.Printf("%s 33 %s\n", args[0], args[2])
			return nil
		}
		val.laps[val.lapsCounter].enterTime = time

		raceStat.participants[competitor] = val

	case 11:
		log.Printf("%s The competitor(%s) can`t continue: %s\n", args[0], args[2], args[3])
	default:
		log.Printf("Wrong event ID\n")
		return fmt.Errorf("wrong event ID")
	}

	return nil
}

// scanning and calling handler of one request
func proc(raceStat *biathlon, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lst := strings.SplitN(line, " ", 4)

		// fmt.Println(len(lst), lst)
		handle(raceStat, lst...)
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}
