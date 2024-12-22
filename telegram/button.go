package telegram

import (
	"fmt"
	"html"
	"sort"
	"strconv"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/fb2k"
	"github.com/c0re100/RadioBot/helper"
	"github.com/c0re100/RadioBot/utils"
	tdlib "github.com/c0re100/gotdlib/client"
)

func createSongListButton(offset int) [][]*tdlib.InlineKeyboardButton {
	var songKb [][]*tdlib.InlineKeyboardButton

	mutex.Lock()
	for i := offset; i < offset+config.GetRowLimit(); i++ {
		if songList[i] == nil {
			continue
		}
		num := strconv.Itoa(i + 1)
		idx := strconv.Itoa(i)
		songKb = append(songKb, []*tdlib.InlineKeyboardButton{
			helper.NewInlineKeyboardCallbackColumn(num, "select_song:"+idx),
		})

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
			"<b>%v</b>. <code>%v</code>", rList, keys[i]+1, html.EscapeString(format))
	}

	return rList
}

func finalizeButton(songKb [][]*tdlib.InlineKeyboardButton, offset int, noBtn bool, sType int) *tdlib.ReplyMarkupInlineKeyboard {
	cbTag := btnTag(sType)

	if noBtn || len(songKb) < config.GetRowLimit() && offset == 0 {

	} else if offset == 0 {
		songKb = append(songKb, []*tdlib.InlineKeyboardButton{
			helper.NewInlineKeyboardCallbackColumn("Next page", cbTag+strconv.Itoa(offset+config.GetRowLimit())),
		})
	} else if len(songKb) < config.GetRowLimit() {
		songKb = append(songKb, []*tdlib.InlineKeyboardButton{
			helper.NewInlineKeyboardCallbackColumn("Previous page", cbTag+strconv.Itoa(offset-config.GetRowLimit())),
		})
	} else {
		songKb = append(songKb, []*tdlib.InlineKeyboardButton{
			helper.NewInlineKeyboardCallbackColumn("Previous page", cbTag+strconv.Itoa(offset-config.GetRowLimit())),
			helper.NewInlineKeyboardCallbackColumn("Next page", cbTag+strconv.Itoa(offset+config.GetRowLimit())),
		})
	}
	return helper.NewReplyMarkupInlineKeyboard(songKb)
}

func btnTag(sType int) string {
	cbTag := "page:"

	switch sType {
	case 0:
		cbTag = "all:"
	case 1:
		cbTag = "artist:"
	case 2:
		cbTag = "track:"
	case 3:
		cbTag = "album:"
	}
	return cbTag
}

func sendButtonMessage(chatID, msgID int64) {
	var format *tdlib.FormattedText
	rList := createResultList(songList, 0)
	if chatID < 0 {
		text := fmt.Sprintf("Which song do you want to play?"+
			"\n\n"+
			"<b>Use Private Chat to request a song WHEN you exceeded rate-limit.</b>\n"+
			"%v", rList)
		format, _ = tdlib.ParseTextEntities(helper.NewHTMLText(text))
	} else {
		text := fmt.Sprintf("Which song do you want to play?\n"+
			"%v", rList)
		format, _ = tdlib.ParseTextEntities(helper.NewHTMLText(text))
	}
	msg := helper.NewEnititesMessage(chatID, msgID, format)
	songKb := createSongListButton(0)
	kb := finalizeButton(songKb, 0, false, 0)
	msg.ReplyMarkup = kb
	_, _ = bot.SendMessage(msg)
}

func editButtonMessage(chatID, msgID int64, queryID tdlib.JsonInt64, offset int, dontCount bool) {
	if canSelectPage(chatID, queryID, dontCount) {
		var format *tdlib.FormattedText
		rList := createResultList(songList, offset)
		if chatID < 0 {
			text := fmt.Sprintf("Which song do you want to play?"+
				"\n\n"+
				"<b>Use Private Chat to request a song WHEN you exceeded rate-limit.</b>\n"+
				"%v", rList)
			format, _ = tdlib.ParseTextEntities(helper.NewHTMLText(text))
		} else {
			text := fmt.Sprintf("Which song do you want to play?\n"+
				"%v", rList)
			format, _ = tdlib.ParseTextEntities(helper.NewHTMLText(text))
		}
		songKb := createSongListButton(offset)
		kb := finalizeButton(songKb, offset, false, 0)
		msg := helper.NewEditMessageText(chatID, msgID, kb, format)
		_, _ = bot.EditMessageText(msg)
	}
}

func selectSongMessage(userID int64, queryID tdlib.JsonInt64, idx int) {
	if songList[idx] == nil {
		_, _ = bot.AnswerCallbackQuery(helper.NewAnswerCallbackQuery(queryID, "This song is not available...", false, "", 180))
	} else if len(GetQueue()) >= config.GetQueueLimit() {
		_, _ = bot.AnswerCallbackQuery(helper.NewAnswerCallbackQuery(queryID, "Too many song in request song list now...\nPlease try again later~", false, "", 180))
	} else {
		if utils.ContainsInt(GetRecent(), idx) {
			_, _ = bot.AnswerCallbackQuery(helper.NewAnswerCallbackQuery(queryID, "Song was recently played!", false, "", 180))
		} else if utils.ContainsInt(GetQueue(), idx) {
			_, _ = bot.AnswerCallbackQuery(helper.NewAnswerCallbackQuery(queryID, "Song was recently requested!", false, "", 180))
		} else {
			if ok, sec := canReqSong(userID); !ok {
				_, _ = bot.AnswerCallbackQuery(helper.NewAnswerCallbackQuery(queryID, fmt.Sprintf("You're already requested recently, Please try again in %v seconds...", sec), false, "", 10))
				return
			}

			fb2k.PushQueue(idx)
			choice := fmt.Sprintf("Your choice: %v - %v | Song queue: %v", songList[idx].Artist, songList[idx].Track, len(GetQueue()))
			_, _ = bot.AnswerCallbackQuery(helper.NewAnswerCallbackQuery(queryID, choice, false, "", 180))
		}
	}
}

func createTypeButton() *tdlib.ReplyMarkupInlineKeyboard {
	kb := [][]*tdlib.InlineKeyboardButton{
		{
			helper.NewInlineKeyboardCallbackColumn("Artist", "select_artist"),
			helper.NewInlineKeyboardCallbackColumn("Track", "select_track"),
			helper.NewInlineKeyboardCallbackColumn("Album", "select_album"),
		},
		{
			helper.NewInlineKeyboardCallbackColumn("Why not both?", "select_all"),
		},
	}
	return helper.NewReplyMarkupInlineKeyboard(kb)
}
