package config

func GetStatus() Status {
	return status
}

func (s Status) GetCurrent() string {
	return s.Current
}

func SetCurrentSong(value string) {
	status.Current = value
}

func (s Status) GetQueue() []int {
	return s.Queue
}

func SetQueueSong(value []int) {
	status.Queue = value
}

func (s Status) GetRecent() []int {
	return s.Recent
}

func SetRecentSong(value []int) {
	status.Recent = value
}
