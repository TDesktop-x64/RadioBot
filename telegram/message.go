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
			switch {
			case command == "/request":
				sendButtonMessage(chatId, msgId)
			case command == "/current":
				getCurrentPlaying(chatId, msgId)
			case command == "/skip":
				startVote(chatId, msgId, int32(senderId))
			case command == "/search" || command == "/nom":
				nominate(chatId, msgId, int32(senderId), CommandArgument(msgText))
			case command == "/play":
				playerControl(chatId, int32(senderId), 0)
			case command == "/stop":
				playerControl(chatId, int32(senderId), 1)
			case command == "/pause":
				playerControl(chatId, int32(senderId), 2)
			case command == "/random":
				playerControl(chatId, int32(senderId), 3)
			case command == "/reload":
				reload(chatId, msgId, int32(senderId))
			case command == "/loadptcps":
				loadParticipants(chatId, int32(senderId))
			}
		}(newMsg)
	}
}
