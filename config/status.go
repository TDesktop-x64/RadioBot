package config

// GetStatus get current status
func GetStatus() Status {
	return status
}

// GetCurrent get current song
func (s Status) GetCurrent() string {
	return s.Current
}

// SetCurrentSong set current song
func SetCurrentSong(value string) {
	status.Current = value
}

// GetQueue get queue song list
func (s Status) GetQueue() []int {
	return s.Queue
}

// SetQueueSong set queue song list
func SetQueueSong(value []int) {
	status.Queue = value
}

// GetRecent get recent song list
func (s Status) GetRecent() []int {
	return s.Recent
}

// SetRecentSong set recent song list
func SetRecentSong(value []int) {
	status.Recent = value
}
