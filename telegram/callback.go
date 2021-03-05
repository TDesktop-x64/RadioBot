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
			queryID := updateMsg.Id
			chatID := updateMsg.ChatId
			userID := updateMsg.SenderUserId
			msgID := updateMsg.MessageId
			data := string(updateMsg.Payload.(*tdlib.CallbackQueryPayloadData).Data)

			m, err := bot.GetMessage(chatID, msgID)
			if err != nil {
				return
			}

			page := strings.Split(data, "page:")
			selIdx := strings.Split(data, "select_song:")
			result := strings.Split(data, "result:")

			switch {
			case data == "vote_skip":
				setUserVote(chatID, msgID, userID, queryID)
			case len(selIdx) == 2:
				idx, _ := strconv.Atoi(selIdx[1])
				selectSongMessage(userID, queryID, idx)
			case len(page) == 2:
				offset, _ := strconv.Atoi(page[1])
				editButtonMessage(chatID, msgID, queryID, offset)
			case len(result) == 2:
				offset, _ := strconv.Atoi(result[1])
				editCustomButtonMessage(chatID, m, queryID, offset)
			}
		}(newMsg)
	}
}
