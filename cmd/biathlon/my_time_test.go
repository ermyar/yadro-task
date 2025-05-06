package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyTimeParse(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    string
		expected myTime
		err      bool
	}{
		{
			name:  "wrong format",
			input: "23:45:59:999",
			err:   true,
		},
		{name: "missed",
			input: "22::00",
			err:   true,
		},
		{
			name:  "wrong hours",
			input: "25:13:24.777",
			err:   true,
		},
		{
			name:  "wrong minutes",
			input: "23:99:59:999",
			err:   true,
		},
		{
			name:  "wrong seconds",
			input: "23:45:-59:999",
			err:   true,
		},
		{
			name:     "typical time",
			input:    "09:30:00.184",
			expected: myTime(34200184),
		},
		{
			name:     "cutted time",
			input:    "23:59:01",
			expected: myTime(86341000),
		},
	} {
		t.Run(tc.name,
			func(t *testing.T) {
				val, err := parseTime(tc.input)
				if tc.err {
					assert.Error(t, err, "must be error!")
					return
				}
				assert.NoError(t, err, "mustn't be an error!")
				assert.Equal(t, tc.expected, val)
			})
	}
}

func TestTimeGetString(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    []myTime
		expected []string
	}{
		{
			name:  "random generated",
			input: []myTime{45256789, 12345678, 59999999, 0, 86399999, 7260450, 31415926, 11223344, 7654321, 557799},
			expected: []string{"12:34:16.789", "03:25:45.678", "16:39:59.999", "00:00:00.000", "23:59:59.999",
				"02:01:00.450", "08:43:35.926", "03:07:03.344", "02:07:34.321", "00:09:17.799"},
		},
	} {
		t.Run(tc.name,
			func(t *testing.T) {
				for i, time := range tc.input {
					assert.Equal(t, tc.expected[i], time.getString())
				}
			})
	}
}
