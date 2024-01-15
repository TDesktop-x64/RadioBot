package helper

import tdlib "github.com/c0re100/gotdlib/client"

func NewMessageSenderUser(userId int64) *tdlib.MessageSenderUser {
	return &tdlib.MessageSenderUser{UserId: userId}
}

func NewMessageSenderChat(chatId int64) *tdlib.MessageSenderChat {
	return &tdlib.MessageSenderChat{ChatId: chatId}
}
