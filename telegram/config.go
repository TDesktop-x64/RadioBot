package telegram

import (
	"fmt"
	"log"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/go-tdlib"
)

func boolToEmoji(b bool) string {
	if b {
		return "✅"
	}
	return "❎"
}

func configFormattedText() (*tdlib.FormattedText, error) {
	currentConf := fmt.Sprintf("<u>Current config</u>\n")
	currentConf += fmt.Sprintf("<b>Chat ID</b>: <code>%v</code>\n", config.GetChatID())
	currentConf += fmt.Sprintf("<b>Chat Username</b>: <code>%v</code>\n", config.GetChatUsername())
	currentConf += fmt.Sprintf("<b>Playlist ID</b>: <code>%v</code>\n", config.GetPlaylistID())
	currentConf += fmt.Sprintf("<b>Song count</b>: <code>%v</code>\n", len(songList))
	currentConf += fmt.Sprintf("\n")
	currentConf += fmt.Sprintf("<u>Limit</u>\n")
	currentConf += fmt.Sprintf("<b>Select page of rate limit for Chat</b>: <code>%v</code>\n", config.GetChatSelectLimit())
	currentConf += fmt.Sprintf("<b>Select page of rate limit for Private Chat</b>: <code>%v</code>\n", config.GetPrivateChatSelectLimit())
	currentConf += fmt.Sprintf("<b>Number of rows</b>: <code>%v</code>\n", config.GetRowLimit())
	currentConf += fmt.Sprintf("<b>Max queue songs</b>: <code>%v</code>\n", config.GetQueueLimit())
	currentConf += fmt.Sprintf("<b>Max recent songs</b>: <code>%v</code>\n", config.GetRecentLimit())
	currentConf += fmt.Sprintf("<b>Request a songs per minute limit</b>: <code>%v</code>\n", config.GetReqSongLimit())
	currentConf += fmt.Sprintf("\n")
	currentConf += fmt.Sprintf("<u>Vote</u>\n")
	currentConf += fmt.Sprintf("<b>Enable</b>: %v\n", boolToEmoji(config.IsVoteEnabled()))
	if config.IsVoteEnabled() {
		currentConf += fmt.Sprintf("<b>Vote time</b>: <code>%v</code>\n", config.GetVoteTime())
		currentConf += fmt.Sprintf("<b>Update vote status each n seconds</b>: <code>%v</code>\n", config.GetUpdateTime())
		currentConf += fmt.Sprintf("<b>Lock the vote n seconds after vote ended</b>: <code>%v</code>\n", config.GetReleaseTime())
		currentConf += fmt.Sprintf("<b>Success percentage</b>: <code>%v%%</code>\n", config.GetSuccessRate())
		currentConf += fmt.Sprintf("<b>Only participants which are in a voice chat can vote</b>: %v\n", boolToEmoji(config.IsPtcpsOnly()))
		currentConf += fmt.Sprintf("<b>Only users which are in the group can vote</b>: %v\n", boolToEmoji(config.IsJoinNeeded()))
	}

	format, err := bot.ParseTextEntities(currentConf, tdlib.NewTextParseModeHTML())
	if err != nil {
		return nil, err
	}

	return format, nil
}

func configButton() *tdlib.ReplyMarkupInlineKeyboard {
	kb := [][]tdlib.InlineKeyboardButton{
		{
			*tdlib.NewInlineKeyboardButton("Refresh", tdlib.NewInlineKeyboardButtonTypeCallback([]byte("refresh_config"))),
		},
		{
			*tdlib.NewInlineKeyboardButton("Vote setting", tdlib.NewInlineKeyboardButtonTypeCallback([]byte("------------------------"))),
		},
		{
			*tdlib.NewInlineKeyboardButton("Enable", tdlib.NewInlineKeyboardButtonTypeCallback([]byte("------------------------"))),
			*tdlib.NewInlineKeyboardButton(boolToEmoji(config.IsVoteEnabled()), tdlib.NewInlineKeyboardButtonTypeCallback([]byte("vote_change"))),
		},
	}

	if config.IsVoteEnabled() {
		kb = append(kb, [][]tdlib.InlineKeyboardButton{
			{
				*tdlib.NewInlineKeyboardButton("Participants only", tdlib.NewInlineKeyboardButtonTypeCallback([]byte("------------------------"))),
				*tdlib.NewInlineKeyboardButton(boolToEmoji(config.IsPtcpsOnly()), tdlib.NewInlineKeyboardButtonTypeCallback([]byte("ptcp_change"))),
			},
			{
				*tdlib.NewInlineKeyboardButton("User join needed", tdlib.NewInlineKeyboardButtonTypeCallback([]byte("------------------------"))),
				*tdlib.NewInlineKeyboardButton(boolToEmoji(config.IsJoinNeeded()), tdlib.NewInlineKeyboardButtonTypeCallback([]byte("join_change"))),
			},
		}...)
	}

	kb = append(kb, []tdlib.InlineKeyboardButton{
		*tdlib.NewInlineKeyboardButton("Reload Config", tdlib.NewInlineKeyboardButtonTypeCallback([]byte("reload_config"))),
		*tdlib.NewInlineKeyboardButton("Reload Playlist", tdlib.NewInlineKeyboardButtonTypeCallback([]byte("reload_playlist"))),
	})

	return tdlib.NewReplyMarkupInlineKeyboard(kb)
}

func configMenu(chatID, msgID int64, userID int32, refresh bool) {
	if !isAdmin(config.GetChatID(), userID) {
		return
	}

	format, err := configFormattedText()
	if err != nil {
		log.Println(err)
		return
	}

	configKb := configButton()
	text := tdlib.NewInputMessageText(format, false, false)
	if refresh {
		bot.EditMessageText(chatID, msgID, configKb, text)
	} else {
		bot.SendMessage(chatID, 0, msgID, tdlib.NewMessageSendOptions(false, true, nil), configKb, text)
	}
}
