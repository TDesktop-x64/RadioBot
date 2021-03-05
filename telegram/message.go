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
			chatId := updateMsg.Message.ChatId
			msgId := updateMsg.Message.Id
			senderId := GetSenderId(updateMsg.Message.Sender)
			var msgText string
			var msgEnt []tdlib.TextEntity

			switch updateMsg.Message.Content.GetMessageContentEnum() {
			case "messageText":
				msgText = updateMsg.Message.Content.(*tdlib.MessageText).Text.Text
				msgEnt = updateMsg.Message.Content.(*tdlib.MessageText).Text.Entities
			case "messageChatJoinByLink":
				bot.DeleteMessages(chatId, []int64{msgId}, true)
			case "messageChatAddMembers", "messageChatDeleteMember":
				bot.DeleteMessages(chatId, []int64{msgId}, true)
			}

			command := CheckCommand(msgText, msgEnt)
			switch command {
			case "/request":
				sendButtonMessage(chatId, msgId)
			case "/current":
				getCurrentPlaying(chatId, msgId)
			case "/skip":
				startVote(chatId, msgId, int32(senderId))
			case "/search", "/nom":
				nominate(chatId, msgId, int32(senderId), CommandArgument(msgText))
			case "/play":
				playerControl(chatId, int32(senderId), 0)
			case "/stop":
				playerControl(chatId, int32(senderId), 1)
			case "/pause":
				playerControl(chatId, int32(senderId), 2)
			case "/random":
				playerControl(chatId, int32(senderId), 3)
			case "/reload":
				reload(chatId, msgId, int32(senderId))
			case "/loadptcps":
				loadParticipants(chatId, int32(senderId))
			}
		}(newMsg)
	}
}
