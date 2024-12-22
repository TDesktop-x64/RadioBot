package helper

import tdlib "github.com/c0re100/gotdlib/client"

func NewInlineKeyboardButton(text string, typ tdlib.InlineKeyboardButtonType) *tdlib.InlineKeyboardButton {
	return &tdlib.InlineKeyboardButton{
		Text: text,
		Type: typ,
	}
}

func NewInlineKeyboardButtonCallback(b []byte) *tdlib.InlineKeyboardButtonTypeCallback {
	return &tdlib.InlineKeyboardButtonTypeCallback{
		Data: b,
	}
}

func NewInlineKeyboardButtonUrl(url string) *tdlib.InlineKeyboardButtonTypeUrl {
	return &tdlib.InlineKeyboardButtonTypeUrl{
		Url: url,
	}
}

func NewInlineKeyboardButtonUser(id int64) *tdlib.InlineKeyboardButtonTypeUser {
	return &tdlib.InlineKeyboardButtonTypeUser{
		UserId: id,
	}
}

func NewInlineKeyboardButtonPassword(b []byte) *tdlib.InlineKeyboardButtonTypeCallbackWithPassword {
	return &tdlib.InlineKeyboardButtonTypeCallbackWithPassword{
		Data: b,
	}
}

func NewInlineKeyboardButtonInline(query string, isCurrentChat bool) *tdlib.InlineKeyboardButtonTypeSwitchInline {
	return &tdlib.InlineKeyboardButtonTypeSwitchInline{
		Query:      query,
		TargetChat: &tdlib.TargetChatCurrent{},
	}
}

func NewInlineKeyboardButtonBuy() *tdlib.InlineKeyboardButtonTypeBuy {
	return &tdlib.InlineKeyboardButtonTypeBuy{}
}

func NewInlineKeyboardButtonLogin(url string, id int64, forwardText string) *tdlib.InlineKeyboardButtonTypeLoginUrl {
	return &tdlib.InlineKeyboardButtonTypeLoginUrl{
		Url:         url,
		Id:          id,
		ForwardText: forwardText,
	}
}

func NewInlineKeyboardButtonGame() *tdlib.InlineKeyboardButtonTypeCallbackGame {
	return &tdlib.InlineKeyboardButtonTypeCallbackGame{}
}

func NewReplyMarkupInlineKeyboard(buttons [][]*tdlib.InlineKeyboardButton) *tdlib.ReplyMarkupInlineKeyboard {
	return &tdlib.ReplyMarkupInlineKeyboard{
		Rows: buttons,
	}
}

func NewAnswerCallbackQuery(callbackQueryId tdlib.JsonInt64, text string, showAlert bool, url string, cacheTime int32) *tdlib.AnswerCallbackQueryRequest {
	return &tdlib.AnswerCallbackQueryRequest{
		CallbackQueryId: callbackQueryId,
		Text:            text,
		ShowAlert:       showAlert,
		Url:             url,
		CacheTime:       cacheTime,
	}
}

func NewInlineKeyboardCallbackRow(text string, callback string) []*tdlib.InlineKeyboardButton {
	return []*tdlib.InlineKeyboardButton{
		{
			Text: text,
			Type: NewInlineKeyboardButtonCallback([]byte(callback)),
		},
	}
}

func NewInlineKeyboardCallbackColumn(text string, callback string) *tdlib.InlineKeyboardButton {
	return &tdlib.InlineKeyboardButton{
		Text: text,
		Type: NewInlineKeyboardButtonCallback([]byte(callback)),
	}
}
