package statistics

type Monitor interface {
	Increment(metric string)
	Close()
}

type Statistics struct {
	monitor Monitor
}

func (s *Statistics) Increment(metric string) {
	s.monitor.Increment(metric)
}

func (s *Statistics) Stop() {
	s.monitor.Close()
}

func (s *Statistics) SetMonitor(monitor Monitor) {
	s.monitor = monitor
}

func NewStatistics() *Statistics {
	return &Statistics{
		monitor: NewStatsHatMonitor(),
	}
}

var statistics = NewStatistics()

func TrackStandingsRequest() {
	statistics.Increment("standings")
}

func TrackLeaguesRequest() {
	statistics.Increment("leagues")
}

func GetApiKey() string {
	return apiKey
}
