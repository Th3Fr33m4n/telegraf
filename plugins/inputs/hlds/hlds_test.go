package hlds

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/influxdata/telegraf/testutil"
)

const testInput = `CPU   In    Out   Uptime  Users   FPS    Players
1.0  0.00  0.00       7     5   98.00       0`

var (
	expectedOutput = statsData{
		1.0, 0.00, 0.00, 7, 5, 98.00, 0,
	}
)

func TestCPUStats(t *testing.T) {
	c := NewHLDSStats()
	var acc testutil.Accumulator
	err := c.gatherServer(&acc, c.Servers[0], requestMock)
	if err != nil {
		t.Error(err)
	}

	if !acc.HasMeasurement("hlds") {
		t.Errorf("acc.HasMeasurement: expected hlds")
	}

	require.Equal(t, "1.2.3.4:1234", acc.Metrics[0].Tags["host"])
	require.Equal(t, "sv1", acc.Metrics[0].Tags["svname"])
	require.Equal(t, expectedOutput.CPU, acc.Metrics[0].Fields["cpu"])
	require.Equal(t, expectedOutput.NetIn, acc.Metrics[0].Fields["net_in"])
	require.Equal(t, expectedOutput.NetOut, acc.Metrics[0].Fields["net_out"])
	require.Equal(t, expectedOutput.UptimeMinutes, acc.Metrics[0].Fields["uptime_minutes"])
	require.Equal(t, expectedOutput.Users, acc.Metrics[0].Fields["users"])
	require.Equal(t, expectedOutput.FPS, acc.Metrics[0].Fields["fps"])
	require.Equal(t, expectedOutput.Players, acc.Metrics[0].Fields["players"])
}

func requestMock(_ string, _ string) (string, error) {
	return testInput, nil
}

func NewHLDSStats() *HLDS {
	return &HLDS{
		Servers: [][]string{
			{"1.2.3.4:1234", "password", "sv1"},
		},
	}
}
