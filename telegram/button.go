package telegram

import (
	"fmt"
	"strconv"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/fb2k"
	"github.com/c0re100/RadioBot/utils"
	"github.com/c0re100/go-tdlib"
)

func createSongListButton(offset int) [][]tdlib.InlineKeyboardButton {
	var songKb [][]tdlib.InlineKeyboardButton

	mutex.Lock()
	for i := offset; i < offset+config.GetRowLimit(); i++ {
		if songList[i] == "" {
			continue
		}
		songKb = append(songKb, []tdlib.InlineKeyboardButton{*tdlib.NewInlineKeyboardButton(songList[i], tdlib.NewInlineKeyboardButtonTypeCallback([]byte("select_song:"+strconv.Itoa(i))))})
	}
	mutex.Unlock()

	return songKb
}

func finalizeButton(songKb [][]tdlib.InlineKeyboardButton, offset int, isSearch bool) *tdlib.ReplyMarkupInlineKeyboard {
	cbTag := "page:"
	if isSearch {
		cbTag = "result:"
	}
	if len(songKb) < 10 && offset == 0 && isSearch {

	} else if offset == 0 {
		songKb = append(songKb, []tdlib.InlineKeyboardButton{
			*tdlib.NewInlineKeyboardButton("Next page", tdlib.NewInlineKeyboardButtonTypeCallback([]byte(cbTag+strconv.Itoa(offset+10)))),
		})
	} else if len(songKb) < 10 {
		songKb = append(songKb, []tdlib.InlineKeyboardButton{
			*tdlib.NewInlineKeyboardButton("Previous page", tdlib.NewInlineKeyboardButtonTypeCallback([]byte(cbTag+strconv.Itoa(offset-10)))),
		})
	} else {
		songKb = append(songKb, []tdlib.InlineKeyboardButton{
			*tdlib.NewInlineKeyboardButton("Previous page", tdlib.NewInlineKeyboardButtonTypeCallback([]byte(cbTag+strconv.Itoa(offset-10)))),
			*tdlib.NewInlineKeyboardButton("Next page", tdlib.NewInlineKeyboardButtonTypeCallback([]byte(cbTag+strconv.Itoa(offset+10)))),
		})
	}
	return tdlib.NewReplyMarkupInlineKeyboard(songKb)
}

func sendButtonMessage(chatID, msgID int64) {
	var format *tdlib.FormattedText
	if chatID < 0 {
		text := "Which song do you want to play?" +
			"\n\n" +
			"<b>Use Private Chat to request a song WHEN you exceeded rate-limit.</b>"
		format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
	} else {
		format = tdlib.NewFormattedText("Which song do you want to play?", nil)
	}
	text := tdlib.NewInputMessageText(format, false, false)
	songKb := createSongListButton(0)
	kb := finalizeButton(songKb, 0, false)
	bot.SendMessage(chatID, 0, msgID, tdlib.NewMessageSendOptions(false, true, nil), kb, text)
}

func editButtonMessage(chatID, msgID int64, queryID tdlib.JSONInt64, offset int) {
	if canSelectPage(chatID, queryID) {
		var format *tdlib.FormattedText
		if chatID < 0 {
			text := "Which song do you want to play?" +
				"\n\n" +
				"<b>Use Private Chat to request a song WHEN you exceeded rate-limit.</b>"
			format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
		} else {
			format = tdlib.NewFormattedText("Which song do you want to play?", nil)
		}
		text := tdlib.NewInputMessageText(format, false, false)
		songKb := createSongListButton(offset)
		kb := finalizeButton(songKb, offset, false)
		bot.EditMessageText(chatID, msgID, kb, text)
	}
}

func selectSongMessage(userID int32, queryID tdlib.JSONInt64, idx int) {
	if ok, sec := canReqSong(userID); !ok {
		bot.AnswerCallbackQuery(queryID, fmt.Sprintf("You're already requested recently, Please try again in %v seconds...", sec), false, "", 59)
		return
	}

	if songList[idx] == "" {
		bot.AnswerCallbackQuery(queryID, "This song is not available...", false, "", 180)
	} else if len(GetQueue()) >= config.GetQueueLimit() {
		bot.AnswerCallbackQuery(queryID, "Too many song in request song list now...\nPlease try again later~", false, "", 180)
	} else {
		if utils.ContainsInt(GetRecent(), idx) {
			bot.AnswerCallbackQuery(queryID, "Song was recently played!", false, "", 180)
		} else if utils.ContainsInt(GetQueue(), idx) {
			bot.AnswerCallbackQuery(queryID, "Song was recently requested!", false, "", 180)
		} else {
			fb2k.PushQueue(idx)
			choice := fmt.Sprintf("Your choice: %v | Song queue: %v", songList[idx], len(GetQueue()))
			bot.AnswerCallbackQuery(queryID, choice, false, "", 180)
		}
	}
}
