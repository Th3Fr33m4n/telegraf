//go:generate ../../../tools/readme_config_includer/generator
package hlds

import (
	_ "embed"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Th3Fr33m4n/udp_rcon"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type statsData struct {
	CPU           float64 `json:"cpu"`
	NetIn         float64 `json:"net_in"`
	NetOut        float64 `json:"net_out"`
	UptimeMinutes float64 `json:"uptime_minutes"`
	Users         float64 `json:"users"`
	FPS           float64 `json:"fps"`
	Players       float64 `json:"players"`
}

type HLDS struct {
	Servers [][]string `toml:"servers"`
}

func (*HLDS) SampleConfig() string {
	return sampleConfig
}

func (s *HLDS) Gather(acc telegraf.Accumulator) error {
	var wg sync.WaitGroup

	// Loop through each server and collect metrics
	for _, server := range s.Servers {
		wg.Add(1)
		go func(ss []string) {
			defer wg.Done()
			acc.AddError(s.gatherServer(acc, ss, requestServer))
		}(server)
	}

	wg.Wait()
	return nil
}

func init() {
	inputs.Add("hlds", func() telegraf.Input {
		return &HLDS{}
	})
}

func (s *HLDS) gatherServer(
	acc telegraf.Accumulator,
	server []string,
	request func(string, string) (string, error),
) error {
	if len(server) != 3 {
		return errors.New("incorrect server config")
	}

	url, rconPw, svID := server[0], server[1], server[2]
	resp, err := request(url, rconPw)
	if err != nil {
		return err
	}

	rows := strings.Split(resp, "\n")
	if len(rows) < 2 {
		return errors.New("bad response")
	}

	fields := strings.Fields(rows[1])
	if len(fields) != 7 {
		return errors.New("bad response")
	}

	cpu, err := strconv.ParseFloat(fields[0], 32)
	if err != nil {
		return err
	}
	netIn, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return err
	}
	netOut, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return err
	}
	uptimeMinutes, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return err
	}
	users, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return err
	}
	fps, err := strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return err
	}
	players, err := strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return err
	}

	now := time.Now()
	stats := statsData{
		CPU:           cpu,
		NetIn:         netIn,
		NetOut:        netOut,
		UptimeMinutes: uptimeMinutes,
		Users:         users,
		FPS:           fps,
		Players:       players,
	}

	tags := map[string]string{
		"host": url,
		"svid": svID,
	}

	var statsMap map[string]interface{}
	marshalled, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshalled, &statsMap)
	if err != nil {
		return err
	}

	acc.AddGauge("hlds", statsMap, tags, now)
	return nil
}

func requestServer(url string, rconPw string) (string, error) {
	remoteConsole, err := rcon.NewRemoteConsole(url, rconPw, true, rcon.UdpConnector{})
	if err != nil {
		return "", err
	}
	defer remoteConsole.Disconnect()

	response, err := remoteConsole.RunCommand("stats", 2048)
	if err != nil {
		return "", err
	}

	return string((*response)[5:]), nil
}
