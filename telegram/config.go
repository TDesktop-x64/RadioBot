package telegram

import (
	"fmt"
	"log"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/helper"
	tdlib "github.com/c0re100/gotdlib/client"
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

	format, err := tdlib.ParseTextEntities(helper.NewHTMLText(currentConf))
	if err != nil {
		return nil, err
	}

	return format, nil
}

func configButton() *tdlib.ReplyMarkupInlineKeyboard {
	kb := [][]*tdlib.InlineKeyboardButton{
		{
			helper.NewInlineKeyboardCallbackColumn("Refresh", "refresh_config"),
		},
		{
			helper.NewInlineKeyboardCallbackColumn("Vote setting", "------------------------"),
		},
		{
			helper.NewInlineKeyboardCallbackColumn("Enable", "------------------------"),
			helper.NewInlineKeyboardCallbackColumn(boolToEmoji(config.IsVoteEnabled()), "vote_change"),
		},
	}

	if config.IsVoteEnabled() {
		kb = append(kb, [][]*tdlib.InlineKeyboardButton{
			{
				helper.NewInlineKeyboardCallbackColumn("Participants only", "------------------------"),
				helper.NewInlineKeyboardCallbackColumn(boolToEmoji(config.IsPtcpsOnly()), "ptcp_change"),
			},
			{
				helper.NewInlineKeyboardCallbackColumn("User join needed", "------------------------"),
				helper.NewInlineKeyboardCallbackColumn(boolToEmoji(config.IsJoinNeeded()), "join_change"),
			},
		}...)
	}

	kb = append(kb, []*tdlib.InlineKeyboardButton{
		helper.NewInlineKeyboardCallbackColumn("Reload Config", "reload_config"),
		helper.NewInlineKeyboardCallbackColumn("Reload Playlist", "reload_playlist"),
	})

	return helper.NewReplyMarkupInlineKeyboard(kb)
}

func configMenu(chatID, msgID int64, userID int64, refresh bool) {
	if !isAdmin(config.GetChatID(), userID) {
		return
	}

	format, err := configFormattedText()
	if err != nil {
		log.Println(err)
		return
	}

	configKb := configButton()
	if refresh {
		msg := helper.NewEditMessageText(chatID, msgID, configKb, format)
		_, _ = bot.EditMessageText(msg)
	} else {
		msg := helper.NewEnititesMessage(chatID, msgID, format)
		msg.ReplyMarkup = configKb
		_, _ = bot.SendMessage(msg)
	}
}
