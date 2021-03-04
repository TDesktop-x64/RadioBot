package config

type Status struct {
	Current string `json:"current"`
	Queue   []int  `json:"queue"`
	Recent  []int  `json:"recent"`
}

type Config struct {
	ApiId        string `json:"api_id"`
	ApiHash      string `json:"api_hash"`
	BotToken     string `json:"bot_token"`
	ChatId       int64  `json:"chat_id"`
	ChatUsername string `json:"chat_username"`
	PinnedMsg    int64  `json:"pinned_message"`
	PlaylistId   string `json:"playlist_id"`
	LimitSetting Limit  `json:"limit"`
	VoteSetting  Vote   `json:"vote"`
	WebSetting   Web    `json:"web"`
}

type Limit struct {
	ChatSelectLimit int `json:"chat_select_limit"`
	PriSelectLimit  int `json:"private_select_limit"`
	RowLimit        int `json:"row_limit"`
	QueueLimit      int `json:"queue_limit"`
	RecentLimit     int `json:"recent_limit"`
	ReqSongPerMin   int `json:"request_song_per_minute"`
}

type Vote struct {
	Enable       bool    `json:"enable"`
	VoteTime     int32   `json:"vote_time"`
	UpdateTime   int32   `json:"update_time"`
	ReleaseTime  int64   `json:"release_time"`
	PctOfSuccess float64 `json:"percent_of_success"`
	PtcpsOnly    bool    `json:"participants_only"`
	UserMustJoin bool    `json:"user_must_join"`
}

type Web struct {
	Enable bool `json:"enable"`
	Port   int  `json:"port"`
}
