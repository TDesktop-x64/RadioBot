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

func canSelectPage(chatID int64, queryID tdlib.JSONInt64, dontCount bool) bool {
	if dontCount {
		return true
	}
	var cID int32

	if chatID == config.GetChatID() {
		cID = -1000
		if pageLimit[cID] == nil {
			pageLimit[cID] = rate.New(config.GetChatSelectLimit(), 1*time.Minute)
		}
	} else if chatID > 0 {
		cID = int32(chatID)
		if pageLimit[cID] == nil {
			pageLimit[cID] = rate.New(config.GetPrivateChatSelectLimit(), 1*time.Minute)
		}
	} else {
		return false
	}

	if ok, dur := pageLimit[cID].Try(); !ok {
		sec := int32(dur.Seconds())
		bot.AnswerCallbackQuery(queryID, fmt.Sprintf("Rate limited! Please try again in %v seconds~", sec), false, "", sec)
		return false
	}
	return true
}

func canReqSong(userID int32) (bool, int) {
	if reqLimit[userID] != nil {
		ok, sec := reqLimit[userID].Try()
		return ok, int(sec.Seconds())
	}
	reqLimit[userID] = rate.New(config.GetReqSongLimit(), 1*time.Minute)
	reqLimit[userID].Try()
	return true, 0
}

func resetRateLimiter() {
	pageLimit = make(map[int32]*rate.RateLimiter)
	reqLimit = make(map[int32]*rate.RateLimiter)
}
