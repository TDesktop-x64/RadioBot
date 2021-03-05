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
		songKb = append(songKb, []tdlib.InlineKeyboardButton{*tdlib.NewInlineKeyboardButton(list[k], tdlib.NewInlineKeyboardButtonTypeCallback([]byte("select_song:"+strconv.Itoa(k))))})
		count++
	}

	return songKb
}

func sendCustomButtonMessage(chatId, msgId int64, list map[int]string) {
	var format *tdlib.FormattedText
	if chatId < 0 {
		text := fmt.Sprintf("Result: %v matches\n"+
			"Which song do you want to play?\n"+
			"\n"+
			"<b>Use Private Chat to request a song WHEN you exceeded rate-limit.</b>", len(list))
		format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
	} else {
		format = tdlib.NewFormattedText(fmt.Sprintf("Result: %v matches\nWhich song do you want to play?", len(list)), nil)
	}
	text := tdlib.NewInputMessageText(format, false, false)
	songKb := createSearchSongListButton(list, 0)

	var kb *tdlib.ReplyMarkupInlineKeyboard
	if len(list) > 0 {
		kb = finalizeButton(songKb, 0, true)
	} else {
		kb = tdlib.NewReplyMarkupInlineKeyboard(songKb)
	}

	bot.SendMessage(chatId, 0, msgId, tdlib.NewMessageSendOptions(false, true, nil), kb, text)
}

func editCustomButtonMessage(chatId int64, m *tdlib.Message, queryId tdlib.JSONInt64, offset int) {
	if canSelectPage(chatId, queryId) {
		var format *tdlib.FormattedText
		if chatId < 0 {
			text := "Which song do you want to play?" +
				"\n\n" +
				"<b>Use Private Chat to request a song WHEN you exceeded rate-limit.</b>"
			format, _ = bot.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
		} else {
			format = tdlib.NewFormattedText("Which song do you want to play?", nil)
		}
		text := tdlib.NewInputMessageText(format, false, false)
		m2, err := bot.GetMessage(chatId, m.ReplyToMessageId)
		if err != nil {
			return
		}
		switch m2.Content.GetMessageContentEnum() {
		case "messageText":
			msgText := m2.Content.(*tdlib.MessageText).Text.Text
			list := searchSong(commandArgument(msgText))
			songKb := createSearchSongListButton(list, offset)
			kb := finalizeButton(songKb, offset, true)
			bot.EditMessageText(chatId, m.Id, kb, text)
		}
	}
}
