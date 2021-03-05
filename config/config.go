package config

func GetConfig() Config {
	return config
}

func GetApiId() string {
	return config.ApiId
}

func GetApiHash() string {
	return config.ApiHash
}

func GetBotToken() string {
	return config.BotToken
}

func GetChatId() int64 {
	return config.ChatId
}

func SetChatId(value int64) {
	config.ChatId = value
}

func GetChatUsername() string {
	return config.ChatUsername
}

func GetPinnedMessage() int64 {
	return config.PinnedMsg << 20
}

func SetPinnedMessage(value int64) {
	config.PinnedMsg = value >> 20
}

func GetBeefWebPort() int {
	return config.BeefWebPort
}

func GetPlaylistId() string {
	return config.PlaylistId
}

func GetChatSelectLimit() int {
	return config.LimitSetting.ChatSelectLimit
}

func GetPrivateChatSelectLimit() int {
	return config.LimitSetting.PriSelectLimit
}

func GetRowLimit() int {
	return config.LimitSetting.RowLimit
}

func GetQueueLimit() int {
	return config.LimitSetting.QueueLimit
}

func GetRecentLimit() int {
	return config.LimitSetting.RecentLimit
}

func GetReqSongLimit() int {
	return config.LimitSetting.ReqSongPerMin
}

func IsVoteEnabled() bool {
	return config.VoteSetting.Enable
}

func GetSuccessRate() float64 {
	return config.VoteSetting.PctOfSuccess
}

func IsPtcpsOnly() bool {
	return config.VoteSetting.PtcpsOnly
}

func GetVoteTime() int32 {
	return config.VoteSetting.VoteTime
}

func GetReleaseTime() int64 {
	return config.VoteSetting.ReleaseTime
}

func GetUpdateTime() int32 {
	return config.VoteSetting.UpdateTime
}

func IsJoinNeeded() bool {
	return config.VoteSetting.UserMustJoin
}

func IsWebEnabled() bool {
	return config.WebSetting.Enable
}

func GetWebPort() int {
	return config.WebSetting.Port
}
