package telegram

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/go-tdlib"
)

func searchSong(text string) map[int]string {
	var list = make(map[int]string)
	for i, s := range songList {
		if strings.Contains(strings.ToLower(s), strings.ToLower(text)) {
			list[i] = s
		}
	}
	return list
}

func createSearchSongListButton(list map[int]string, offset int) [][]tdlib.InlineKeyboardButton {
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
	for _, k := range keys[offset:] {
		if count >= config.GetRowLimit() {
			break
		}
		if list[k] == "" {
			continue
		}
		idx := strconv.Itoa(k)
		songKb = append(songKb, []tdlib.InlineKeyboardButton{*tdlib.NewInlineKeyboardButton(idx, tdlib.NewInlineKeyboardButtonTypeCallback([]byte("select_song:"+idx)))})
		count++
	}

	return songKb
}

func sendCustomButtonMessage(chatID, msgID int64, list map[int]string) {
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
	if len(list) > 0 {
		kb = finalizeButton(songKb, 0, true)
	} else {
		kb = tdlib.NewReplyMarkupInlineKeyboard(songKb)
	}

	bot.SendMessage(chatID, 0, msgID, tdlib.NewMessageSendOptions(false, true, nil), kb, text)
}

func editCustomButtonMessage(chatID int64, m *tdlib.Message, queryID tdlib.JSONInt64, offset int) {
	if canSelectPage(chatID, queryID) {
		m2, err := bot.GetMessage(chatID, m.ReplyToMessageId)
		if err != nil {
			return
		}
		switch m2.Content.GetMessageContentEnum() {
		case "messageText":
			msgText := m2.Content.(*tdlib.MessageText).Text.Text
			list := searchSong(commandArgument(msgText))
			rList := createResultList(list, offset)

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
			kb := finalizeButton(songKb, offset, true)
			bot.EditMessageText(chatID, m.Id, kb, text)
		}
	}
}

func nominate(chatID, msgID int64, userID int32, arg string) {
	if arg == "" {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Track name or Artist name is empty.", nil), true, false)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		return
	}

	list := searchSong(arg)
	if len(list) > 0 {
		sendCustomButtonMessage(chatID, msgID, list)
	} else {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("No result.", nil), true, false)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
	}
}
