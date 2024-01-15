package telegram

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/helper"
	"github.com/c0re100/RadioBot/utils"
	tdlib "github.com/c0re100/gotdlib/client"
)

func searchAll(text string) map[int]*songInfo {
	var list = make(map[int]*songInfo)
	for i, s := range songList {
		if strings.Contains(strings.ToLower(s.Artist), strings.ToLower(text)) {
			list[i] = s
			continue
		}
		if strings.Contains(strings.ToLower(s.Track), strings.ToLower(text)) {
			list[i] = s
			continue
		}
		if strings.Contains(strings.ToLower(s.Album), strings.ToLower(text)) {
			list[i] = s
		}
	}
	return list
}

func searchArtist(text string) map[int]*songInfo {
	var list = make(map[int]*songInfo)
	for i, s := range songList {
		if strings.Contains(strings.ToLower(s.Artist), strings.ToLower(text)) {
			list[i] = s
		}
	}
	return list
}

func searchTrack(text string) map[int]*songInfo {
	var list = make(map[int]*songInfo)
	for i, s := range songList {
		if strings.Contains(strings.ToLower(s.Track), strings.ToLower(text)) {
			list[i] = s
		}
	}
	return list
}

func searchAlbum(text string) map[int]*songInfo {
	var list = make(map[int]*songInfo)
	for i, s := range songList {
		if strings.Contains(strings.ToLower(s.Album), strings.ToLower(text)) {
			list[i] = s
		}
	}
	return list
}

func createSearchSongListButton(list map[int]*songInfo, offset int) [][]*tdlib.InlineKeyboardButton {
	var songKb [][]*tdlib.InlineKeyboardButton

	if offset > len(list) {
		return songKb
	}

	keys := make([]int, 0, len(list))
	for k := range list {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	count := 0
	for _, i := range keys[offset:] {
		if count >= config.GetRowLimit() {
			break
		}
		if list[i] == nil {
			continue
		}
		num := strconv.Itoa(i + 1)
		idx := strconv.Itoa(i)
		songKb = append(songKb, []*tdlib.InlineKeyboardButton{
			helper.NewInlineKeyboardCallbackColumn(num, "select_song:"+idx),
		})
		count++
	}

	return songKb
}

func sendCustomButtonMessage(chatID, msgID int64, list map[int]*songInfo, sType int) {
	var format *tdlib.FormattedText
	rList := createResultList(list, 0)
	if chatID < 0 {
		text := fmt.Sprintf("Result: %v matches\n"+
			"Which song do you want to play?\n"+
			"\n"+
			"<b>Use Private Chat to request a song WHEN you exceeded rate-limit.</b>\n"+
			"%v", len(list), rList)
		format, _ = tdlib.ParseTextEntities(helper.NewHTMLText(text))
	} else {
		text := fmt.Sprintf("Result: %v matches\n"+
			"Which song do you want to play?\n"+
			"%v", len(list), rList)
		format, _ = tdlib.ParseTextEntities(helper.NewHTMLText(text))
	}
	songKb := createSearchSongListButton(list, 0)

	var kb *tdlib.ReplyMarkupInlineKeyboard
	if len(list) > config.GetRowLimit() {
		kb = finalizeButton(songKb, 0, false, sType)
	} else if len(list) <= config.GetRowLimit() {
		kb = finalizeButton(songKb, 0, true, sType)
	} else {
		kb = helper.NewReplyMarkupInlineKeyboard(songKb)
	}

	msg := helper.NewEditMessageText(chatID, msgID, kb, format)
	_, _ = bot.EditMessageText(msg)
}

func editCustomButtonMessage(chatID int64, m *tdlib.Message, queryID tdlib.JsonInt64, offset int, sType int) {
	if canSelectPage(chatID, queryID, false) {
		m2, err := bot.GetMessage(helper.NewGetMessage(chatID, utils.GetReplyMessageId(m.ReplyTo)))
		if err != nil {
			if sType == 0 {
				switch m.Content.MessageContentType() {
				case "messageText":
					editButtonMessage(chatID, m.Id, queryID, offset, true)
				}
			} else {
				_, _ = bot.AnswerCallbackQuery(helper.NewAnswerCallbackQuery(queryID, "Please search again~", false, "", 0))
			}
			return
		}
		switch m2.Content.MessageContentType() {
		case "messageText":
			msgText := m2.Content.(*tdlib.MessageText).Text.Text

			var list map[int]*songInfo
			var rList string

			list = createSearchList(sType, list, msgText)
			rList = createResultList(list, offset)

			var format *tdlib.FormattedText
			if chatID < 0 {
				text := fmt.Sprintf("Result: %v matches\n"+
					"Which song do you want to play?\n"+
					"\n"+
					"<b>Use Private Chat to request a song WHEN you exceeded rate-limit.</b>\n"+
					"%v", len(list), rList)
				format, _ = tdlib.ParseTextEntities(helper.NewHTMLText(text))
			} else {
				text := fmt.Sprintf("Result: %v matches\n"+
					"Which song do you want to play?\n"+
					"%v", len(list), rList)
				format, _ = tdlib.ParseTextEntities(helper.NewHTMLText(text))
			}

			songKb := createSearchSongListButton(list, offset)
			kb := createResultKeyboard(sType, songKb, offset)
			msg := helper.NewEditMessageText(chatID, m.Id, kb, format)
			_, _ = bot.EditMessageText(msg)
		}
	}
}

func createResultKeyboard(sType int, songKb [][]*tdlib.InlineKeyboardButton, offset int) *tdlib.ReplyMarkupInlineKeyboard {
	var kb *tdlib.ReplyMarkupInlineKeyboard
	switch sType {
	case 0:
		kb = finalizeButton(songKb, offset, false, 0)
	case 1:
		kb = finalizeButton(songKb, offset, false, 1)
	case 2:
		kb = finalizeButton(songKb, offset, false, 2)
	case 3:
		kb = finalizeButton(songKb, offset, false, 3)
	}
	return kb
}

func createSearchList(sType int, list map[int]*songInfo, msgText string) map[int]*songInfo {
	arg := tdlib.CommandArgument(msgText)
	switch sType {
	case 0:
		list = searchAll(arg)
	case 1:
		list = searchArtist(arg)
	case 2:
		list = searchTrack(arg)
	case 3:
		list = searchAlbum(arg)
	}
	return list
}

func nominateType(chatID, msgID int64, userID int64, arg string) {
	if arg == "" {
		msg := helper.NewSimpleMessage(chatID, 0, msgID, "Track/Artist/Album is empty.")
		_, _ = bot.SendMessage(msg)
		return
	}

	msg := helper.NewSimpleMessage(chatID, 0, msgID, "Select type to search")
	_, _ = bot.SendMessage(msg)
}

func valueIsEmpty(chatID, msgID int64, arg string) bool {
	if arg == "" {
		msg := helper.NewSimpleMessage(chatID, 0, msgID, "Value is empty.")
		_, _ = bot.SendMessage(msg)
		return true
	}
	return false
}

func nominate(chatID, msgID int64, userID int64, arg string) {
	if valueIsEmpty(chatID, msgID, arg) {
		return
	}

	list := searchAll(arg)
	if len(list) > 0 {
		sendCustomButtonMessage(chatID, msgID, list, 0)
	} else {
		msg := helper.NewEditMessageText(chatID, msgID, nil, helper.NewFormattedText("No result.", nil))
		_, _ = bot.EditMessageText(msg)
	}
}

func nominateArtist(chatID, msgID int64, userID int64, arg string) {
	if valueIsEmpty(chatID, msgID, arg) {
		return
	}

	list := searchArtist(arg)
	if len(list) > 0 {
		sendCustomButtonMessage(chatID, msgID, list, 1)
	} else {
		msg := helper.NewEditMessageText(chatID, msgID, nil, helper.NewFormattedText("No result.", nil))
		_, _ = bot.EditMessageText(msg)
	}
}

func nominateTrack(chatID, msgID int64, userID int64, arg string) {
	if valueIsEmpty(chatID, msgID, arg) {
		return
	}

	list := searchTrack(arg)
	if len(list) > 0 {
		sendCustomButtonMessage(chatID, msgID, list, 2)
	} else {
		msg := helper.NewEditMessageText(chatID, msgID, nil, helper.NewFormattedText("No result.", nil))
		_, _ = bot.EditMessageText(msg)
	}
}

func nominateAlbum(chatID, msgID int64, userID int64, arg string) {
	if valueIsEmpty(chatID, msgID, arg) {
		return
	}

	list := searchAlbum(arg)
	if len(list) > 0 {
		sendCustomButtonMessage(chatID, msgID, list, 3)
	} else {
		msg := helper.NewEditMessageText(chatID, msgID, nil, helper.NewFormattedText("No result.", nil))
		_, _ = bot.EditMessageText(msg)
	}
}
