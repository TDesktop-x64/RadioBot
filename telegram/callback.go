package telegram

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/c0re100/go-tdlib"
)

func callbackQuery() {
	fmt.Println("[Music] New Callback Receiver")
	eventFilter := func(msg *tdlib.TdMessage) bool {
		return true
	}
	receiver := bot.AddEventReceiver(&tdlib.UpdateNewCallbackQuery{}, eventFilter, 1000)
	for newMsg := range receiver.Chan {
		go func(newMsg tdlib.TdMessage) {
			updateMsg := (newMsg).(*tdlib.UpdateNewCallbackQuery)
			queryId := updateMsg.Id
			chatId := updateMsg.ChatId
			userId := updateMsg.SenderUserId
			msgId := updateMsg.MessageId
			data := string(updateMsg.Payload.(*tdlib.CallbackQueryPayloadData).Data)

			m, err := bot.GetMessage(chatId, msgId)
			if err != nil {
				return
			}

			page := strings.Split(data, "page:")
			selIdx := strings.Split(data, "select_song:")
			result := strings.Split(data, "result:")

			switch {
			case data == "vote_skip":
				setUserVote(chatId, msgId, userId, queryId)
			case len(selIdx) == 2:
				idx, _ := strconv.Atoi(selIdx[1])
				selectSongMessage(userId, queryId, idx)
			case len(page) == 2:
				offset, _ := strconv.Atoi(page[1])
				editButtonMessage(chatId, msgId, queryId, offset)
			case len(result) == 2:
				offset, _ := strconv.Atoi(result[1])
				editCustomButtonMessage(chatId, m, queryId, offset)
			}
		}(newMsg)
	}
}
