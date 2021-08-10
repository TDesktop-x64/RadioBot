package telegram

import (
	"fmt"
	"sort"
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
		if songList[i] == nil {
			continue
		}
		num := strconv.Itoa(i + 1)
		idx := strconv.Itoa(i)
		songKb = append(songKb, []tdlib.InlineKeyboardButton{*tdlib.NewInlineKeyboardButton(num, tdlib.NewInlineKeyboardButtonTypeCallback([]byte("select_song:"+idx)))})
	}
	mutex.Unlock()

	return songKb
}

func createResultList(list map[int]*songInfo, offset int) string {
	var rList string
	var keys []int

	for k := range list {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for i := offset; i < offset+config.GetRowLimit(); i++ {
		if keys == nil {
			break
		}
		if len(keys) == i {
			break
		}
		format := list[keys[i]].Artist + " - " + list[keys[i]].Track
		rList = fmt.Sprintf("%v\n"+
			"<b>%v</b>. <code>%v</code>", rList, keys[i]+1, format)
	}

	return rList
}

func finalizeButton(songKb [][]tdlib.InlineKeyboardButton, offset int, isSearch, noBtn, isAlbum bool) *tdlib.ReplyMarkupInlineKeyboard {
	cbTag := "page:"
	if isAlbum {
		cbTag = "album:"
	} else if isSearch {
		cbTag = "result:"
	}
	if noBtn || len(songKb) < config.GetRowLimit() && offset == 0 && isSearch {

	} else if offset == 0 {
		songKb = append(songKb, []tdlib.InlineKeyboardButton{
			*tdlib.NewInlineKeyboardButton("Next page", tdlib.NewInlineKeyboardButtonTypeCallback([]byte(cbTag+strconv.Itoa(offset+config.GetRowLimit())))),
		})
	} else if len(songKb) < config.GetRowLimit() {
		songKb = append(songKb, []tdlib.InlineKeyboardButton{
			*tdlib.NewInlineKeyboardButton("Previous page", tdlib.NewInlineKeyboardButtonTypeCallback([]byte(cbTag+strconv.Itoa(offset-config.GetRowLimit())))),
		})
	} else {
		songKb = append(songKb, []tdlib.InlineKeyboardButton{
			*tdlib.NewInlineKeyboardButton("Previous page", tdlib.NewInlineKeyboardButtonTypeCallback([]byte(cbTag+strconv.Itoa(offset-config.GetRowLimit())))),
			*tdlib.NewInlineKeyboardButton("Next page", tdlib.NewInlineKeyboardButtonTypeCallback([]byte(cbTag+strconv.Itoa(offset+config.GetRowLimit())))),
		})
	}
	return tdlib.NewReplyMarkupInlineKeyboard(songKb)
}

func sendButtonMessage(chatID, msgID int64) {
	var format *tdlib.FormattedText
	rList := createResultList(songList, 0)
	if chatID < 0 {
		text := fmt.Sprintf("Which song do you want to play?"+
			"\n\n"+
			"<b>Use Private Chat to request a song WHEN you exceeded rate-limit.</b>\n"+
			"%v", rList)
		format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
	} else {
		text := fmt.Sprintf("Which song do you want to play?\n"+
			"%v", rList)
		format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
	}
	text := tdlib.NewInputMessageText(format, false, false)
	songKb := createSongListButton(0)
	kb := finalizeButton(songKb, 0, false, false, false)
	bot.SendMessage(chatID, 0, msgID, tdlib.NewMessageSendOptions(false, true, nil), kb, text)
}

func editButtonMessage(chatID, msgID int64, queryID tdlib.JSONInt64, offset int) {
	if canSelectPage(chatID, queryID) {
		var format *tdlib.FormattedText
		rList := createResultList(songList, offset)
		if chatID < 0 {
			text := fmt.Sprintf("Which song do you want to play?"+
				"\n\n"+
				"<b>Use Private Chat to request a song WHEN you exceeded rate-limit.</b>\n"+
				"%v", rList)
			format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
		} else {
			text := fmt.Sprintf("Which song do you want to play?\n"+
				"%v", rList)
			format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
		}
		text := tdlib.NewInputMessageText(format, false, false)
		songKb := createSongListButton(offset)
		kb := finalizeButton(songKb, offset, false, false, false)
		bot.EditMessageText(chatID, msgID, kb, text)
	}
}

func selectSongMessage(userID int32, queryID tdlib.JSONInt64, idx int) {
	if songList[idx] == nil {
		bot.AnswerCallbackQuery(queryID, "This song is not available...", false, "", 180)
	} else if len(GetQueue()) >= config.GetQueueLimit() {
		bot.AnswerCallbackQuery(queryID, "Too many song in request song list now...\nPlease try again later~", false, "", 180)
	} else {
		if utils.ContainsInt(GetRecent(), idx) {
			bot.AnswerCallbackQuery(queryID, "Song was recently played!", false, "", 180)
		} else if utils.ContainsInt(GetQueue(), idx) {
			bot.AnswerCallbackQuery(queryID, "Song was recently requested!", false, "", 180)
		} else {
			if ok, sec := canReqSong(userID); !ok {
				bot.AnswerCallbackQuery(queryID, fmt.Sprintf("You're already requested recently, Please try again in %v seconds...", sec), false, "", 10)
				return
			}

			fb2k.PushQueue(idx)
			choice := fmt.Sprintf("Your choice: %v | Song queue: %v", songList[idx], len(GetQueue()))
			bot.AnswerCallbackQuery(queryID, choice, false, "", 180)
		}
	}
}

func createTypeButton() *tdlib.ReplyMarkupInlineKeyboard {
	kb := [][]tdlib.InlineKeyboardButton{
		{
			*tdlib.NewInlineKeyboardButton("Track/Artist", tdlib.NewInlineKeyboardButtonTypeCallback([]byte("select_all"))),
			*tdlib.NewInlineKeyboardButton("Album", tdlib.NewInlineKeyboardButtonTypeCallback([]byte("select_album"))),
		},
	}
	return tdlib.NewReplyMarkupInlineKeyboard(kb)
}
