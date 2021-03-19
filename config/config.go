package config

import "log"

// GetConfig get current config
func GetConfig() Config {
	return config
}

// GetAPIID get API ID
func GetAPIID() string {
	return config.ApiId
}

// GetAPIHash get API Hash
func GetAPIHash() string {
	return config.ApiHash
}

// GetBotToken get bot token
func GetBotToken() string {
	return config.BotToken
}

// GetChatID get chat ID
func GetChatID() int64 {
	return config.ChatId
}

// SetChatID update chat ID
func SetChatID(value int64) {
	config.ChatId = value
}

// GetChatUsername get chat username
func GetChatUsername() string {
	return config.ChatUsername
}

// GetPinnedMessage get pinned message ID
func GetPinnedMessage() int64 {
	return config.PinnedMsg << 20
}

// SetPinnedMessage update pinned message ID
func SetPinnedMessage(value int64) {
	config.PinnedMsg = value >> 20
}

// GetBeefWebPort get Beefweb port
func GetBeefWebPort() int {
	return config.BeefwebPort
}

// GetPlaylistID get playlist ID
func GetPlaylistID() string {
	return config.PlaylistId
}

// SetChatSelectLimit set select limit of group chat
func SetChatSelectLimit(value int) {
	config.LimitSetting.ChatSelectLimit = value
}

// GetChatSelectLimit get select limit of group chat
func GetChatSelectLimit() int {
	return config.LimitSetting.ChatSelectLimit
}

// SetPrivateChatSelectLimit set select limit of private chat
func SetPrivateChatSelectLimit(value int) {
	config.LimitSetting.PriSelectLimit = value
}

// GetPrivateChatSelectLimit get select limit of private chat
func GetPrivateChatSelectLimit() int {
	return config.LimitSetting.PriSelectLimit
}

// SetRowLimit set row limit
func SetRowLimit(value int) {
	config.LimitSetting.RowLimit = value
}

// GetRowLimit get row limit
func GetRowLimit() int {
	return config.LimitSetting.RowLimit
}

// SetQueueLimit set queue song limit
func SetQueueLimit(value int) {
	config.LimitSetting.QueueLimit = value
}

// GetQueueLimit get queue song limit
func GetQueueLimit() int {
	return config.LimitSetting.QueueLimit
}

// SetRecentLimit set recent song limit
func SetRecentLimit(value int) {
	config.LimitSetting.RecentLimit = value
}

// GetRecentLimit get recent song limit
func GetRecentLimit() int {
	return config.LimitSetting.RecentLimit
}

// SetReqSongLimit set request song limit
func SetReqSongLimit(value int) {
	config.LimitSetting.ReqSongPerMin = value
}

// GetReqSongLimit get request song limit
func GetReqSongLimit() int {
	return config.LimitSetting.ReqSongPerMin
}

// IsVoteEnabled return true if vote is enabled
func IsVoteEnabled() bool {
	return config.VoteSetting.Enable
}

// SetVoteEnable set vote on or off
func SetVoteEnable(value bool) {
	config.VoteSetting.Enable = value
}

// SetSuccessRate set vote success rate
func SetSuccessRate(value float64) {
	config.VoteSetting.PctOfSuccess = value
}

// GetSuccessRate get vote success rate
func GetSuccessRate() float64 {
	return config.VoteSetting.PctOfSuccess
}

// SetPtcpEnable update participants only
func SetPtcpEnable(value bool) {
	config.VoteSetting.PtcpsOnly = value
}

// IsPtcpsOnly return true if only participants which are in a voice chat can vote
func IsPtcpsOnly() bool {
	return config.VoteSetting.PtcpsOnly
}

// SetVoteTime update the vote time
func SetVoteTime(value int32) {
	config.VoteSetting.VoteTime = value
}

// GetVoteTime get vote time
func GetVoteTime() int32 {
	return config.VoteSetting.VoteTime
}

// SetReleaseTime set lock the vote seconds after vote ended
func SetReleaseTime(value int64) {
	config.VoteSetting.ReleaseTime = value
}

// GetReleaseTime get lock the vote seconds after vote ended
func GetReleaseTime() int64 {
	return config.VoteSetting.ReleaseTime
}

// SetUpdateTime update the vote update time
func SetUpdateTime(value int32) {
	config.VoteSetting.UpdateTime = value
}

// GetUpdateTime get vote update time
func GetUpdateTime() int32 {
	return config.VoteSetting.UpdateTime
}

// IsJoinNeeded return true if only users which are in the group can vote
func IsJoinNeeded() bool {
	return config.VoteSetting.UserMustJoin
}

// SetJoinEnable update user must join
func SetJoinEnable(value bool) {
	config.VoteSetting.UserMustJoin = value
}

// IsWebEnabled return true if userbot mode is disabled
func IsWebEnabled() bool {
	return config.WebSetting.Enable
}

// GetWebPort get web port
func GetWebPort() int {
	return config.WebSetting.Port
}

func compareUpdateVoteTime() {
	if GetUpdateTime() > GetVoteTime() {
		SetUpdateTime(5)
		log.Println("'update_time' is greater than 'vote_time' is not allowed.\n" +
			"Applying default value(5s) to 'update_time'.")
	}
}

func checkVoteTimeIsTooSmall() {
	if GetVoteTime() < 5 {
		SetVoteTime(5)
		log.Println("'vote_time' is smaller than 5s is not allowed.\n" +
			"Value increased to 5s")
	}
}

func checkUpdateTimeIsTooSmall() {
	if GetUpdateTime() < 5 {
		SetUpdateTime(5)
		log.Println("'update_time' is smaller than 5s is not allowed.\n" +
			"Value increased to 5s")
	}
}
