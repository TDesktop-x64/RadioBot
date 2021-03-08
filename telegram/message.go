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
			senderID := getSenderID(updateMsg.Message.Sender)
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
			case "/request":
				sendButtonMessage(chatID, msgID)
			case "/current":
				getCurrentPlaying(chatID, msgID)
			case "/skip":
				startVote(chatID, msgID, int32(senderID))
			case "/search", "/nom":
				nominate(chatID, msgID, int32(senderID), commandArgument(msgText))
			case "/queue":
				checkQueueSong(chatID, msgID)
			case "/play":
				playerControl(chatID, int32(senderID), 0)
			case "/stop":
				playerControl(chatID, int32(senderID), 1)
			case "/pause":
				playerControl(chatID, int32(senderID), 2)
			case "/random":
				playerControl(chatID, int32(senderID), 3)
			case "/reload":
				reload(chatID, msgID, int32(senderID))
			case "/loadptcps":
				loadParticipants(chatID, int32(senderID))
			}
		}(newMsg)
	}
}
