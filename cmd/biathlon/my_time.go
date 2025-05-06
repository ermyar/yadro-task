package main

import (
	"fmt"
	"strconv"
	"strings"
)

// provide Time type in one day
// (hoping that all timing are processing in one day)
type myTime uint64

const (
	miliSec  myTime = 1000
	miliMin  myTime = miliSec * 60
	miliHour myTime = miliMin * 60
)

// returns myTime data from `str`
// which match `HH:MM:SS.sss`
func parseTime(str string) (myTime, error) {
	var ans myTime = 0
	parts := strings.Split(str, ":")

	if len(parts) != 3 {
		return 0, fmt.Errorf("wrong format of TIME")
	}

	if hours, err := strconv.Atoi(parts[0]); err == nil && 0 <= hours && hours < 24 {
		ans += miliHour * myTime(hours)
	} else {
		return 0, fmt.Errorf("error parsing hours")
	}

	if minutes, err := strconv.Atoi(parts[1]); err == nil && 0 <= minutes && minutes < 60 {
		ans += miliMin * myTime(minutes)
	} else {
		return 0, fmt.Errorf("error parsing minutes")
	}

	secFormat := strings.Split(parts[2], ".")

	if seconds, err := strconv.Atoi(secFormat[0]); err == nil && 0 <= seconds && seconds < 60 {
		ans += miliSec * myTime(seconds)
	} else {
		return 0, fmt.Errorf("error parsing seconds")
	}

	if len(secFormat) == 2 {
		if miliSeconds, err := strconv.Atoi(secFormat[1]); err == nil && 0 <= miliSeconds && miliSeconds < 1000 {
			ans += myTime(miliSeconds)
		} else {
			return 0, fmt.Errorf("error parsing seconds")
		}
	}

	return ans, nil
}

// returning time in format `HH:MM:SS.sss`
func (time myTime) getString() (str string) {
	h := time / miliHour
	time %= miliHour

	m := time / miliMin
	time %= miliMin

	s := time / miliSec
	time %= miliSec

	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, time)
}

func (time myTime) sub(other myTime) myTime {
	return time - other
}

func (time myTime) add(interval myTime) myTime {
	return time + interval
}
