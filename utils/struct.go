package utils

type Event struct {
	Player struct {
		ActiveItem struct {
			Columns       []string `json:"columns"`
			Duration      float64  `json:"duration"`
			Index         int      `json:"index"`
			PlaylistID    string   `json:"playlistId"`
			PlaylistIndex int      `json:"playlistIndex"`
			Position      float64  `json:"position"`
		} `json:"activeItem"`
		Info struct {
			Name          string `json:"name"`
			PluginVersion string `json:"pluginVersion"`
			Title         string `json:"title"`
			Version       string `json:"version"`
		} `json:"info"`
		PlaybackMode  int      `json:"playbackMode"`
		PlaybackModes []string `json:"playbackModes"`
		PlaybackState string   `json:"playbackState"`
		Volume        struct {
			IsMuted bool    `json:"isMuted"`
			Max     float64 `json:"max"`
			Min     float64 `json:"min"`
			Type    string  `json:"type"`
			Value   float64 `json:"value"`
		} `json:"volume"`
	} `json:"player"`
}
