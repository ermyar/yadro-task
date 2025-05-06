package integration

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		setup()
		defer teardown()

		return m.Run()
	}())
}

// contains testCase description
type testCase struct {
	name        string
	err         bool
	jsonName    string
	dataSet     string
	expectedAns string
}

func TestBiathlon(t *testing.T) {
	for _, tc := range []testCase{
		{
			name:     "wrong json",
			jsonName: "configs/sample.json",
			dataSet:  "testdata/sample",
			err:      true,
		},
		{
			name:     "wrong dataSet",
			jsonName: "../../configs/sample.json",
			dataSet:  "go.mod",
			err:      true,
		},
		{
			name:        "sample",
			jsonName:    "../../configs/sample.json",
			dataSet:     "testdata/sample",
			expectedAns: "testdata/expected/sample",
		},
		{
			name:        "mixed",
			jsonName:    "../../configs/sample.json",
			dataSet:     "testdata/mixed",
			expectedAns: "testdata/expected/mixed",
		},
		{
			name:        "events",
			jsonName:    "../../configs/config.json",
			dataSet:     "testdata/events",
			expectedAns: "testdata/expected/events",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := []string{"--config", tc.jsonName, "--events", tc.dataSet}
			cmd := exec.Command("./biathlon", args...)
			cmd.Stderr = os.Stderr

			output, err := cmd.Output()

			if tc.err {
				require.Error(t, err)
				_, ok := err.(*exec.ExitError)
				require.True(t, ok)
			} else {
				require.NoError(t, err)
				expected, err := os.ReadFile(tc.expectedAns)
				require.NoError(t, err)
				require.Equal(t, string(expected), string(output))
			}
		})
	}
}

func setup() {
	cmd := exec.Command("go", "build", "../../cmd/biathlon")
	if err := cmd.Run(); err != nil {
		panic("can't build project")
	}
}

func teardown() {
	cmd := exec.Command("rm", "biathlon")
	if err := cmd.Run(); err != nil {
		panic("can't clean")
	}
}
