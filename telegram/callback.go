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

			m2, err2 := bot.GetMessage(chatID, m.ReplyToMessageId)
			if err2 != nil {
				if data == "select_all" || data == "select_album" {
					return
				}
			}

			page := strings.Split(data, "page:")
			selIdx := strings.Split(data, "select_song:")
			result := strings.Split(data, "result:")
			album := strings.Split(data, "album:")

			switch {
			case data == "vote_skip":
				setUserVote(chatID, msgID, userID, queryID)
			case data == "refresh_config":
				configMenu(chatID, msgID, userID, true)
			case data == "reload_config":
				reloadConfig(queryID, userID)
			case data == "reload_playlist":
				reloadPlaylist(queryID, userID)
			case data == "vote_change":
				voteOptionControl(chatID, msgID, userID, 0)
			case data == "ptcp_change":
				voteOptionControl(chatID, msgID, userID, 1)
			case data == "join_change":
				voteOptionControl(chatID, msgID, userID, 2)
			case data == "select_all":
				if m2.Content.GetMessageContentEnum() == "messageText" {
					nominateTrack(chatID, msgID, userID, commandArgument(m2.Content.(*tdlib.MessageText).Text.Text))
				}
			case data == "select_album":
				if m2.Content.GetMessageContentEnum() == "messageText" {
					nominateAlbum(chatID, msgID, userID, commandArgument(m2.Content.(*tdlib.MessageText).Text.Text))
				}
			case len(selIdx) == 2:
				idx, _ := strconv.Atoi(selIdx[1])
				selectSongMessage(userID, queryID, idx)
			case len(page) == 2:
				offset, _ := strconv.Atoi(page[1])
				editButtonMessage(chatID, msgID, queryID, offset)
			case len(result) == 2:
				offset, _ := strconv.Atoi(result[1])
				editCustomButtonMessage(chatID, m, queryID, offset, false)
			case len(album) == 2:
				offset, _ := strconv.Atoi(album[1])
				editCustomButtonMessage(chatID, m, queryID, offset, true)
			}
		}(newMsg)
	}
}
