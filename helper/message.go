package helper

import tdlib "github.com/c0re100/gotdlib/client"

func NewDeleteMessages(chatId int64, messageIds []int64) *tdlib.DeleteMessagesRequest {
	return &tdlib.DeleteMessagesRequest{ChatId: chatId, MessageIds: messageIds, Revoke: true}
}

func NewGetMessage(chatId int64, messageId int64) *tdlib.GetMessageRequest {
	return &tdlib.GetMessageRequest{ChatId: chatId, MessageId: messageId}
}

func NewGetMessages(chatId int64, messageIds []int64) *tdlib.GetMessagesRequest {
	return &tdlib.GetMessagesRequest{ChatId: chatId, MessageIds: messageIds}
}

func NewRemotePhoto(chatId, threadId, msgId int64, fileId string, caption *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessagePhoto{
			Photo: &tdlib.InputFileRemote{
				Id: fileId,
			},
			Caption: caption,
		},
	}
}

func NewRemoteAnimation(chatId, threadId, msgId int64, fileId string, caption *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageAnimation{
			Animation: &tdlib.InputFileRemote{
				Id: fileId,
			},
			Caption: caption,
		},
	}
}

func NewRemoteSticker(chatId, threadId, msgId int64, fileId string) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageSticker{
			Sticker: &tdlib.InputFileRemote{
				Id: fileId,
			},
		},
	}
}

func NewRemoteVideo(chatId, threadId, msgId int64, fileId string, caption *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageVideo{
			Video: &tdlib.InputFileRemote{
				Id: fileId,
			},
			Caption: caption,
		},
	}
}

func NewRemoteDocument(chatId, threadId, msgId int64, fileId string, caption *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageDocument{
			Document: &tdlib.InputFileRemote{
				Id: fileId,
			},
			Caption: caption,
		},
	}
}

func NewRemoteAudio(chatId, threadId, msgId int64, fileId string, caption *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageAudio{
			Audio: &tdlib.InputFileRemote{
				Id: fileId,
			},
			Caption: caption,
		},
	}
}

func NewLocalPhoto(chatId, threadId, msgId int64, path string, caption *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessagePhoto{
			Photo: &tdlib.InputFileLocal{
				Path: path,
			},
			Caption: caption,
		},
	}
}

func NewLocalAnimation(chatId, threadId, msgId int64, path string, caption *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageAnimation{
			Animation: &tdlib.InputFileLocal{
				Path: path,
			},
			Caption: caption,
		},
	}
}

func NewLocalSticker(chatId, threadId, msgId int64, path string) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageSticker{
			Sticker: &tdlib.InputFileLocal{
				Path: path,
			},
		},
	}
}

func NewLocalVideo(chatId, threadId, msgId int64, path string, caption *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageVideo{
			Video: &tdlib.InputFileLocal{
				Path: path,
			},
			Caption: caption,
		},
	}
}

func NewLocalDocument(chatId, threadId, msgId int64, path string, caption *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageDocument{
			Document: &tdlib.InputFileLocal{
				Path: path,
			},
			Caption: caption,
		},
	}
}

func NewLocalAudio(chatId, threadId, msgId int64, path string, caption *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageAudio{
			Audio: &tdlib.InputFileLocal{
				Path: path,
			},
			Caption: caption,
		},
	}
}

func NewFormattedText(text string, entities []*tdlib.TextEntity) *tdlib.FormattedText {
	return &tdlib.FormattedText{
		Text:     text,
		Entities: entities,
	}
}

func NewSimpleMessage(chatId, threadId, msgId int64, text string) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageText{
			Text: &tdlib.FormattedText{
				Text: text,
			},
		},
	}
}

func NewEnititesMessage(chatId, threadId, msgId int64, format *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageText{
			Text: format,
		},
	}
}

func NewInputMessageContent(chatId, threadId, msgId int64, input tdlib.InputMessageContent) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: input,
	}
}

func NewHTMLText(text string) *tdlib.ParseTextEntitiesRequest {
	return &tdlib.ParseTextEntitiesRequest{
		Text:      text,
		ParseMode: &tdlib.TextParseModeHTML{},
	}
}

func NewMarkDownText(text string) *tdlib.ParseTextEntitiesRequest {
	return &tdlib.ParseTextEntitiesRequest{
		Text:      text,
		ParseMode: &tdlib.TextParseModeMarkdown{},
	}
}

func NewReplyMarkupMessage(chatId, threadId, msgId int64, markup tdlib.ReplyMarkup, text string) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		ReplyMarkup: markup,
		InputMessageContent: &tdlib.InputMessageText{
			Text: &tdlib.FormattedText{
				Text: text,
			},
		},
	}
}

func NewPhotoMessage(chatId, threadId, msgId int64, path string, caption *tdlib.FormattedText) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessagePhoto{
			Photo: &tdlib.InputFileLocal{
				Path: path,
			},
			Caption: caption,
		},
	}
}

func NewStickerMessage(chatId, threadId, msgId int64, path string) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: threadId,
		ReplyTo: &tdlib.InputMessageReplyToMessage{
			MessageId: msgId,
		},
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		InputMessageContent: &tdlib.InputMessageSticker{
			Sticker: &tdlib.InputFileLocal{
				Path: path,
			},
		},
	}
}

func NewForwardMessages(chatId, fromChatId, threadId int64, msgId []int64, isCopy bool) *tdlib.ForwardMessagesRequest {
	return &tdlib.ForwardMessagesRequest{
		ChatId:          chatId,
		FromChatId:      fromChatId,
		MessageThreadId: threadId,
		MessageIds:      msgId,
		Options: &tdlib.MessageSendOptions{
			FromBackground: true,
		},
		SendCopy: isCopy,
	}
}

func NewEditMessageText(chatId int64, msgId int64, markup tdlib.ReplyMarkup, format *tdlib.FormattedText) *tdlib.EditMessageTextRequest {
	return &tdlib.EditMessageTextRequest{
		ChatId:      chatId,
		MessageId:   msgId,
		ReplyMarkup: markup,
		InputMessageContent: &tdlib.InputMessageText{
			Text: format,
		},
	}
}

func NewEditMessageMedia(chatId int64, msgId int64, markup tdlib.ReplyMarkup, fileId string, caption *tdlib.FormattedText) *tdlib.EditMessageMediaRequest {
	return &tdlib.EditMessageMediaRequest{
		ChatId:      chatId,
		MessageId:   msgId,
		ReplyMarkup: markup,
		InputMessageContent: &tdlib.InputMessagePhoto{
			Photo: &tdlib.InputFileRemote{
				Id: fileId,
			},
			Caption: caption,
		},
	}
}
