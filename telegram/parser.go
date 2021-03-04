package telegram

import (
	"strings"

	"github.com/c0re100/go-tdlib"
)

func CheckCommand(msgText string, entity []tdlib.TextEntity) string {
	if msgText != "" {
		if msgText[0] == '/' {
			if len(entity) >= 1 {
				if entity[0].Type.GetTextEntityTypeEnum() == "textEntityTypeBotCommand" {
					if i := strings.Index(msgText[:entity[0].Length], "@"); i != -1 {
						return msgText[:i]
					}
					return msgText[:entity[0].Length]
				}
			}
			if len(msgText) > 1 {
				if i := strings.Index(msgText, "@"); i != -1 {
					return msgText[:i]
				}
				if i := strings.Index(msgText, " "); i != -1 {
					return msgText[:i]
				}
				return msgText
			}
		}
	}
	return ""
}

func CommandArgument(msgText string) string {
	if msgText[0] == '/' {
		if i := strings.Index(msgText, " "); i != -1 {
			return msgText[i+1:]
		}
	}
	return ""
}

func GetSenderId(sender tdlib.MessageSender) int64 {
	if sender.GetMessageSenderEnum() == "messageSenderUser" {
		return int64(sender.(*tdlib.MessageSenderUser).UserId)
	} else {
		return sender.(*tdlib.MessageSenderChat).ChatId
	}
}
