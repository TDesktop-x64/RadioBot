package telegram

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/go-tdlib"
)

func searchSong(text string) map[int]*songInfo {
	var list = make(map[int]*songInfo)
	for i, s := range songList {
		if strings.Contains(strings.ToLower(s.Artist), strings.ToLower(text)) {
			list[i] = s
			continue
		} else if strings.Contains(strings.ToLower(s.Track), strings.ToLower(text)) {
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

func sendCustomButtonMessage(chatID, msgID int64, list map[int]*songInfo, isAlbum bool) {
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
		kb = finalizeButton(songKb, 0, true, false, isAlbum)
	} else if len(list) <= config.GetRowLimit() {
		kb = finalizeButton(songKb, 0, true, true, isAlbum)
	} else {
		kb = tdlib.NewReplyMarkupInlineKeyboard(songKb)
	}

	bot.EditMessageText(chatID, msgID, kb, text)
}

func editCustomButtonMessage(chatID int64, m *tdlib.Message, queryID tdlib.JSONInt64, offset int, isAlbum bool) {
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

			if isAlbum {
				list = searchAlbum(commandArgument(msgText))
				rList = createResultList(list, offset)
			} else {
				list = searchSong(commandArgument(msgText))
				rList = createResultList(list, offset)
			}

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

			var kb *tdlib.ReplyMarkupInlineKeyboard
			if isAlbum {
				kb = finalizeButton(songKb, offset, true, false, true)
			} else {
				kb = finalizeButton(songKb, offset, true, false, false)
			}

			bot.EditMessageText(chatID, m.Id, kb, text)
		}
	}
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

	list := searchSong(arg)
	if len(list) > 0 {
		sendCustomButtonMessage(chatID, msgID, list, false)
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
		sendCustomButtonMessage(chatID, msgID, list, true)
	} else {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("No result.", nil), true, false)
		bot.EditMessageText(chatID, msgID, nil, msgText)
	}
}
