package config

// Status status JSON struct
type Status struct {
	Current string `json:"current"`
	Queue   []int  `json:"queue"`
	Recent  []int  `json:"recent"`
}

// Config config JSON strust
type Config struct {
	ApiId        string `json:"api_id"`
	ApiHash      string `json:"api_hash"`
	BotToken     string `json:"bot_token"`
	ChatId       int64  `json:"chat_id"`
	ChatUsername string `json:"chat_username"`
	PinnedMsg    int64  `json:"pinned_message"`
	BeefwebPort  int    `json:"beefweb_port"`
	PlaylistId   string `json:"playlist_id"`
	LimitSetting Limit  `json:"limit"`
	VoteSetting  Vote   `json:"vote"`
	WebSetting   Web    `json:"web"`
}

// Limit limit setting JSON struct
type Limit struct {
	ChatSelectLimit int `json:"chat_select_limit"`
	PriSelectLimit  int `json:"private_select_limit"`
	RowLimit        int `json:"row_limit"`
	QueueLimit      int `json:"queue_limit"`
	RecentLimit     int `json:"recent_limit"`
	ReqSongPerMin   int `json:"request_song_per_minute"`
}

// Vote vote setting JSON struct
type Vote struct {
	Enable       bool    `json:"enable"`
	VoteTime     int32   `json:"vote_time"`
	UpdateTime   int32   `json:"update_time"`
	ReleaseTime  int64   `json:"release_time"`
	PctOfSuccess float64 `json:"percent_of_success"`
	PtcpsOnly    bool    `json:"participants_only"`
	UserMustJoin bool    `json:"user_must_join"`
}

// Web web setting JSON struct
type Web struct {
	Enable bool `json:"enable"`
	Port   int  `json:"port"`
}
