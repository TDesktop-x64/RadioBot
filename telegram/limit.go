package telegram

import (
	"fmt"
	"time"

	"github.com/beefsack/go-rate"
	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/go-tdlib"
)

var (
	pageLimit = make(map[int32]*rate.RateLimiter)
	reqLimit  = make(map[int32]*rate.RateLimiter)
)

func canSelectPage(chatId int64, queryId tdlib.JSONInt64) bool {
	var cId int32

	if chatId == config.GetChatId() {
		cId = -1000
		if pageLimit[cId] == nil {
			pageLimit[cId] = rate.New(config.GetChatSelectLimit(), 1*time.Minute)
		}
	} else if chatId > 0 {
		cId = int32(chatId)
		if pageLimit[cId] == nil {
			pageLimit[cId] = rate.New(config.GetPrivateChatSelectLimit(), 1*time.Minute)
		}
	} else {
		return false
	}

	if ok, dur := pageLimit[cId].Try(); !ok {
		sec := int32(dur.Seconds())
		bot.AnswerCallbackQuery(queryId, fmt.Sprintf("Rate limited! Please try again in %v seconds~", sec), false, "", sec)
		return false
	}
	return true
}

func canReqSong(userId int32) (bool, int) {
	if reqLimit[userId] != nil {
		ok, sec := reqLimit[userId].Try()
		return ok, int(sec)
	}
	reqLimit[userId] = rate.New(config.GetReqSongLimit(), 1*time.Minute)
	reqLimit[userId].Try()
	return true, 0
}
