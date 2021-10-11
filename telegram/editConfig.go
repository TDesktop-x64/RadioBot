package telegram

import (
	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/go-tdlib"
	"strconv"
)

func reloadConfig(queryID tdlib.JSONInt64, userID int64) {
	if !isAdmin(config.GetChatID(), userID) {
		return
	}

	if err := config.LoadConfig(); err != nil {
		bot.AnswerCallbackQuery(queryID, err.Error(), false, "", 10)
		return
	}
	resetRateLimiter()
	bot.AnswerCallbackQuery(queryID, "Config reloaded.", false, "", 10)
}

func reloadPlaylist(queryID tdlib.JSONInt64, userID int64) {
	if !isAdmin(config.GetChatID(), userID) {
		return
	}

	if err := savePlaylistIndexAndName(); err != nil {
		bot.AnswerCallbackQuery(queryID, err.Error(), false, "", 10)
		return
	}
	bot.AnswerCallbackQuery(queryID, "Playlist reloaded.", false, "", 10)
}

func optionControl(chatID, msgID int64, userID int64, cs int, arg string) {
	if !isAdmin(config.GetChatID(), userID) {
		return
	}

	if arg == "" {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Command argument is empty.\nFormat: /setting_name <integer>", nil), true, false)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		return
	}

	val, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Command argument must be integer.\nFormat: /setting_name <integer>", nil), true, false)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		return
	}

	var cmd string
	switch cs {
	case 0:
		cmd = "Limit: Select page of rate limit for Group"
		config.SetChatSelectLimit(int(val))
		resetRateLimiter()
	case 1:
		cmd = "Limit: Select page of rate limit for Private Chat"
		config.SetPrivateChatSelectLimit(int(val))
		resetRateLimiter()
	case 2:
		cmd = "Limit: Number of rows"
		config.SetRowLimit(int(val))
	case 3:
		cmd = "Limit: Max queue songs"
		config.SetQueueLimit(int(val))
	case 4:
		cmd = "Limit: Max recent songs"
		config.SetRecentLimit(int(val))
	case 5:
		cmd = "Limit: Request song per minute"
		config.SetReqSongLimit(int(val))
		resetRateLimiter()
	case 6:
		cmd = "Timer: Vote time"
		config.SetVoteTime(int32(val))
	case 7:
		cmd = "Timer: Update time"
		config.SetUpdateTime(int32(val))
	case 8:
		cmd = "Timer: Release time"
		config.SetReleaseTime(val)
	case 9:
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Command argument must be float.\nFormat: /setting_name <float>", nil), true, false)
			bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
			return
		}
		cmd = "Timer: Vote success rate"
		config.SetSuccessRate(val)
	}

	msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText(cmd+" set to "+arg+".", nil), true, false)
	bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
}

func voteOptionControl(chatID, msgID int64, userID int64, cs int) {
	if !isAdmin(config.GetChatID(), userID) {
		return
	}

	switch cs {
	case 0:
		config.SetVoteEnable(!config.IsVoteEnabled())
	case 1:
		config.SetPtcpEnable(!config.IsPtcpsOnly())
	case 2:
		config.SetJoinEnable(!config.IsJoinNeeded())
	default:
		return
	}

	configMenu(chatID, msgID, userID, true)
}
