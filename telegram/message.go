package telegram

import (
	"fmt"

	"github.com/TDesktop-x64/RadioBot/helper"
	"github.com/TDesktop-x64/RadioBot/utils"
	tdlib "github.com/c0re100/gotdlib/client"
)

func newMessages() {
	fmt.Println("[Music] New Message Receiver")
	listener := bot.AddEventReceiver(&tdlib.UpdateNewMessage{}, 100)
	for newMsg := range listener.Updates {
		go func(newMsg tdlib.Type) {
			updateMsg := (newMsg).(*tdlib.UpdateNewMessage)
			chatID := updateMsg.Message.ChatId
			msgID := updateMsg.Message.Id
			senderId := utils.GetSenderId(updateMsg.Message.SenderId)
			var msgText string
			var msgEnt []*tdlib.TextEntity

			switch updateMsg.Message.Content.MessageContentType() {
			case "messageText":
				msgText = updateMsg.Message.Content.(*tdlib.MessageText).Text.Text
				msgEnt = updateMsg.Message.Content.(*tdlib.MessageText).Text.Entities
			case "messageChatJoinByLink":
				_, _ = bot.DeleteMessages(helper.NewDeleteMessages(chatID, []int64{msgID}))
			case "messageChatAddMembers", "messageChatDeleteMember":
				_, _ = bot.DeleteMessages(helper.NewDeleteMessages(chatID, []int64{msgID}))
			}

			command := tdlib.CheckCommand(msgText, msgEnt)
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
				nominateType(chatID, msgID, senderId, tdlib.CommandArgument(msgText))
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
				optionControl(chatID, msgID, senderId, 0, tdlib.CommandArgument(msgText))
			case "/private_select_limit":
				optionControl(chatID, msgID, senderId, 1, tdlib.CommandArgument(msgText))
			case "/row_limit":
				optionControl(chatID, msgID, senderId, 2, tdlib.CommandArgument(msgText))
			case "/queue_limit":
				optionControl(chatID, msgID, senderId, 3, tdlib.CommandArgument(msgText))
			case "/recent_limit":
				optionControl(chatID, msgID, senderId, 4, tdlib.CommandArgument(msgText))
			case "/request_song_per_minute":
				optionControl(chatID, msgID, senderId, 5, tdlib.CommandArgument(msgText))
			case "/vote_time":
				optionControl(chatID, msgID, senderId, 6, tdlib.CommandArgument(msgText))
			case "/update_time":
				optionControl(chatID, msgID, senderId, 7, tdlib.CommandArgument(msgText))
			case "/release_time":
				optionControl(chatID, msgID, senderId, 8, tdlib.CommandArgument(msgText))
			case "/percent_of_success":
				optionControl(chatID, msgID, senderId, 9, tdlib.CommandArgument(msgText))
			case "/loadptcps":
				loadParticipants(chatID, senderId)
			}
		}(newMsg)
	}
}
