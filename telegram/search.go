package telegram

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/go-tdlib"
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

func createSearchSongListButton(list map[int]*songInfo, offset int) [][]tdlib.InlineKeyboardButton {
	var songKb [][]tdlib.InlineKeyboardButton

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
		songKb = append(songKb, []tdlib.InlineKeyboardButton{*tdlib.NewInlineKeyboardButton(num, tdlib.NewInlineKeyboardButtonTypeCallback([]byte("select_song:"+idx)))})
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
		format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
	} else {
		text := fmt.Sprintf("Result: %v matches\n"+
			"Which song do you want to play?\n"+
			"%v", len(list), rList)
		format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
	}
	text := tdlib.NewInputMessageText(format, false, false)
	songKb := createSearchSongListButton(list, 0)

	var kb *tdlib.ReplyMarkupInlineKeyboard
	if len(list) > config.GetRowLimit() {
		kb = finalizeButton(songKb, 0, false, sType)
	} else if len(list) <= config.GetRowLimit() {
		kb = finalizeButton(songKb, 0, true, sType)
	} else {
		kb = tdlib.NewReplyMarkupInlineKeyboard(songKb)
	}

	bot.EditMessageText(chatID, msgID, kb, text)
}

func editCustomButtonMessage(chatID int64, m *tdlib.Message, queryID tdlib.JSONInt64, offset int, sType int) {
	if canSelectPage(chatID, queryID) {
		m2, err := bot.GetMessage(chatID, m.ReplyToMessageId)
		if err != nil {
			return
		}
		switch m2.Content.GetMessageContentEnum() {
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
				format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
			} else {
				text := fmt.Sprintf("Result: %v matches\n"+
					"Which song do you want to play?\n"+
					"%v", len(list), rList)
				format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
			}

			text := tdlib.NewInputMessageText(format, false, false)
			songKb := createSearchSongListButton(list, offset)
			kb := createResultKeyboard(sType, songKb, offset)
			bot.EditMessageText(chatID, m.Id, kb, text)
		}
	}
}

func createResultKeyboard(sType int, songKb [][]tdlib.InlineKeyboardButton, offset int) *tdlib.ReplyMarkupInlineKeyboard {
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
	switch sType {
	case 0:
		list = searchAll(commandArgument(msgText))
	case 1:
		list = searchArtist(commandArgument(msgText))
	case 2:
		list = searchTrack(commandArgument(msgText))
	case 3:
		list = searchAlbum(commandArgument(msgText))
	}
	return list
}

func nominateType(chatID, msgID int64, userID int32, arg string) {
	if arg == "" {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Track/Artist/Album is empty.", nil), true, false)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		return
	}

	msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Select type to search", nil), true, false)
	bot.SendMessage(chatID, 0, msgID, nil, createTypeButton(), msgText)
}

func valueIsEmpty(chatID, msgID int64, arg string) bool {
	if arg == "" {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Value is empty.", nil), true, false)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		return true
	}
	return false
}

func nominate(chatID, msgID int64, userID int32, arg string) {
	if valueIsEmpty(chatID, msgID, arg) {
		return
	}

	list := searchAll(arg)
	if len(list) > 0 {
		sendCustomButtonMessage(chatID, msgID, list, 0)
	} else {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("No result.", nil), true, false)
		bot.EditMessageText(chatID, msgID, nil, msgText)
	}
}

func nominateArtist(chatID, msgID int64, userID int32, arg string) {
	if valueIsEmpty(chatID, msgID, arg) {
		return
	}

	list := searchArtist(arg)
	if len(list) > 0 {
		sendCustomButtonMessage(chatID, msgID, list, 1)
	} else {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("No result.", nil), true, false)
		bot.EditMessageText(chatID, msgID, nil, msgText)
	}
}

func nominateTrack(chatID, msgID int64, userID int32, arg string) {
	if valueIsEmpty(chatID, msgID, arg) {
		return
	}

	list := searchTrack(arg)
	if len(list) > 0 {
		sendCustomButtonMessage(chatID, msgID, list, 2)
	} else {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("No result.", nil), true, false)
		bot.EditMessageText(chatID, msgID, nil, msgText)
	}
}

func nominateAlbum(chatID, msgID int64, userID int32, arg string) {
	if valueIsEmpty(chatID, msgID, arg) {
		return
	}

	list := searchAlbum(arg)
	if len(list) > 0 {
		sendCustomButtonMessage(chatID, msgID, list, 3)
	} else {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("No result.", nil), true, false)
		bot.EditMessageText(chatID, msgID, nil, msgText)
	}
}
