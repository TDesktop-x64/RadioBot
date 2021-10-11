package telegram

import (
	"fmt"

	"github.com/c0re100/go-tdlib"
)

func newMessages() {
	fmt.Println("[Music] New Message Receiver")
	eventFilter := func(msg *tdlib.TdMessage) bool {
		return true
	}

	receiver := bot.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 100)
	for newMsg := range receiver.Chan {
		go func(newMsg tdlib.TdMessage) {
			updateMsg := (newMsg).(*tdlib.UpdateNewMessage)
			chatID := updateMsg.Message.ChatId
			msgID := updateMsg.Message.Id
			senderId := getsenderId(updateMsg.Message.Sender)
			var msgText string
			var msgEnt []tdlib.TextEntity

			switch updateMsg.Message.Content.GetMessageContentEnum() {
			case "messageText":
				msgText = updateMsg.Message.Content.(*tdlib.MessageText).Text.Text
				msgEnt = updateMsg.Message.Content.(*tdlib.MessageText).Text.Entities
			case "messageChatJoinByLink":
				bot.DeleteMessages(chatID, []int64{msgID}, true)
			case "messageChatAddMembers", "messageChatDeleteMember":
				bot.DeleteMessages(chatID, []int64{msgID}, true)
			}

			command := checkCommand(msgText, msgEnt)
			switch command {
			case "/config":
				configMenu(chatID, msgID, senderId, false)
			case "/request":
				sendButtonMessage(chatID, msgID)
			case "/current":
				getCurrentPlaying(chatID, msgID)
			case "/skip":
				startVote(chatID, msgID, senderId)
			case "/search", "/nom":
				nominateType(chatID, msgID, senderId, commandArgument(msgText))
			case "/queue":
				checkQueueSong(chatID, msgID)
			case "/latest":
				checkLatestSong(chatID, msgID, len(songList)-30)
			case "/play":
				playerControl(chatID, senderId, 0)
			case "/stop":
				playerControl(chatID, senderId, 1)
			case "/pause":
				playerControl(chatID, senderId, 2)
			case "/random":
				playerControl(chatID, senderId, 3)
			case "/chat_select_limit":
				optionControl(chatID, msgID, senderId, 0, commandArgument(msgText))
			case "/private_select_limit":
				optionControl(chatID, msgID, senderId, 1, commandArgument(msgText))
			case "/row_limit":
				optionControl(chatID, msgID, senderId, 2, commandArgument(msgText))
			case "/queue_limit":
				optionControl(chatID, msgID, senderId, 3, commandArgument(msgText))
			case "/recent_limit":
				optionControl(chatID, msgID, senderId, 4, commandArgument(msgText))
			case "/request_song_per_minute":
				optionControl(chatID, msgID, senderId, 5, commandArgument(msgText))
			case "/vote_time":
				optionControl(chatID, msgID, senderId, 6, commandArgument(msgText))
			case "/update_time":
				optionControl(chatID, msgID, senderId, 7, commandArgument(msgText))
			case "/release_time":
				optionControl(chatID, msgID, senderId, 8, commandArgument(msgText))
			case "/percent_of_success":
				optionControl(chatID, msgID, senderId, 9, commandArgument(msgText))
			case "/loadptcps":
				loadParticipants(chatID, senderId)
			}
		}(newMsg)
	}
}
