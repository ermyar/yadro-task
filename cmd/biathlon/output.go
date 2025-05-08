package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

// printing final report for all competitors
// use special formatting by tabwritter (for goodlooking)
func output(raceStat *biathlon, data *jsonRace) {
	sl := make([]int, 0)
	for k := range raceStat.participants {
		sl = append(sl, k)
	}
	sort.Slice(sl, func(i, j int) bool {
		st1 := raceStat.participants[sl[i]]
		st2 := raceStat.participants[sl[j]]
		return st1.status < st2.status || (st1.status == st2.status && st1.totalTime < st2.totalTime)
	})
	fmt.Println("\nResulting Table:")

	w := tabwriter.Writer{}
	w.Init(os.Stdout, 1, 0, 1, ' ', 0)
	for _, val := range sl {
		fmt.Fprintln(&w, getReport(val, raceStat.participants[val], data))
	}

	w.Flush()
}

func getFormatLap(lapLen int, time myTime) string {
	if time == 0 {
		return "{,}"
	}
	str := fmt.Sprintf("%.6f", float64(miliSec)*float64(lapLen)/float64(time))
	return fmt.Sprintf("{%s, %s}", time.getString(), str[:len(str)-3])
}

// return string contain information about `index` participant
// in required format
func getReport(index int, val stat, data *jsonRace) string {
	var total string
	switch val.status {
	case finished:
		total = val.totalTime.getString()
	case notFinished:
		total = "NotFinished"
	case notStarted:
		total = "NotStarted"
	}

	laps := make([]string, data.Laps)
	for i := range laps {
		laps[i] = getFormatLap(data.LapLen, val.laps[i].totalTime)
	}

	penaltyLap := getFormatLap(val.shotStat.missed*data.PenaltyLen, val.penalty.totalTime)

	return fmt.Sprintf("[%s]\t%d\t[%s]\t%s\t%d/%d", total, index, strings.Join(laps, ", "), penaltyLap, val.shotStat.hitted, val.shotStat.total)
}
