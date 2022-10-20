package statistics

import stathat "github.com/stathat/go"

type StatsHatMonitor struct{}

var (
	apiKey = ""
)

func NewStatsHatMonitor() Monitor {
	return &StatsHatMonitor{}
}

func (stats *StatsHatMonitor) Increment(metric string) {
	stathat.PostEZCount(metric, apiKey, 1)
}

func (stats *StatsHatMonitor) Close() {
}
