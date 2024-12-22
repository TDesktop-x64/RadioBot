package telegram

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/fb2k"
	"github.com/c0re100/RadioBot/helper"
	"github.com/c0re100/RadioBot/utils"
	tdlib "github.com/c0re100/gotdlib/client"
)

type groupStatus struct {
	chatID       int64
	msgID        int64
	vcID         int32
	duartion     int32
	Ptcps        []string
	voteSkip     []int64
	isVoting     bool
	isLoadPtcps  bool
	lastVoteTime int64
}

var (
	grpStatus = &groupStatus{chatID: config.GetChatID()}
)

// GetQueue get queue song list
func GetQueue() []int {
	return config.GetStatus().GetQueue()
}

// GetRecent get recent song list
func GetRecent() []int {
	return config.GetStatus().GetRecent()
}

func getUserIDHash(uID int64) string {
	h := sha1.New()
	h.Write([]byte(strconv.FormatInt(int64(uID), 10)))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

func startVote(chatID, msgID int64, userID int64) {
	if chatID != config.GetChatID() {
		return
	}

	if !config.IsVoteEnabled() {
		msg := helper.NewSimpleMessage(chatID, msgID, "This group is not allowed to vote.")
		_, _ = bot.SendMessage(msg)
		return
	}

	if !config.IsWebEnabled() {
		c, err := userBot.GetChat(&tdlib.GetChatRequest{ChatId: chatID})
		if err != nil {
			log.Println(err)
			return
		}

		if c.VideoChat.GroupCallId == 0 {
			msg := helper.NewSimpleMessage(chatID, msgID, "This group do not have a voice chat.")
			_, _ = bot.SendMessage(msg)
			return
		}
		// Preload all users
		_, _ = userBot.LoadGroupCallParticipants(&tdlib.LoadGroupCallParticipantsRequest{GroupCallId: c.VideoChat.GroupCallId, Limit: 5000})
	}

	hashedID := getUserIDHash(int64(userID))
	if config.IsPtcpsOnly() && !utils.ContainsString(grpStatus.Ptcps, hashedID) {
		msg := helper.NewSimpleMessage(chatID, msgID, "Only users which are in a voice chat can vote!")
		_, _ = bot.SendMessage(msg)
		return
	}

	if grpStatus.isVoting {
		msg := helper.NewSimpleMessage(chatID, msgID, "Vote in progress...")
		_, _ = bot.SendMessage(msg)
		return
	}

	if time.Now().Unix() < grpStatus.lastVoteTime+config.GetReleaseTime() {
		msg := helper.NewSimpleMessage(chatID, msgID, "Skip a song was voted too recently...")
		_, _ = bot.SendMessage(msg)
		return
	}

	voteKb := helper.NewReplyMarkupInlineKeyboard([][]*tdlib.InlineKeyboardButton{
		{
			helper.NewInlineKeyboardCallbackColumn("Yes - 1", "vote_skip"),
		},
	})

	msg := helper.NewSimpleMessage(chatID, msgID, "Skip a song?")
	msg.ReplyMarkup = voteKb
	m, err := bot.SendMessage(msg)
	if err != nil {
		log.Println("Can't send message.")
		return
	}
	grpStatus.isVoting = true
	grpStatus.duartion = config.GetVoteTime()
	grpStatus.msgID = m.Id
	grpStatus.lastVoteTime = time.Now().Unix()
	updateTime := config.GetUpdateTime()

	if !utils.ContainsInt64(grpStatus.voteSkip, userID) {
		grpStatus.voteSkip = append(grpStatus.voteSkip, userID)
	}
	updateVote(chatID, m.Id, false)
	// Wait N seconds
	time.Sleep(time.Duration(updateTime) * time.Second)
	addVoteJob(chatID, m.Id, updateTime)
	if !sch.IsRunning() {
		log.Println("Starting scheduler...")
		startScheduler()
	}
}

func updateVote(chatID, msgID int64, isAuto bool) {
	if isAuto {
		grpStatus.duartion -= config.GetUpdateTime()
	}
	if grpStatus.duartion <= 0 {
		endVote(chatID, msgID)
		return
	}
	voteKb := helper.NewReplyMarkupInlineKeyboard([][]*tdlib.InlineKeyboardButton{
		{
			helper.NewInlineKeyboardCallbackColumn(fmt.Sprintf("Yes - %v", len(grpStatus.voteSkip)), "vote_skip"),
		},
	})

	format := fmt.Sprintf("Skip a song?\n"+
		"Vote count: %v\n"+
		"Vote timeleft: %v second(s)", len(grpStatus.voteSkip), grpStatus.duartion)
	msg := helper.NewEditMessageText(chatID, msgID, voteKb, helper.NewFormattedText(format, nil))
	_, _ = bot.EditMessageText(msg)
}

func resetVote() {
	sch.RemoveByTag("timeleft")
	grpStatus.isLoadPtcps = false
	grpStatus.isVoting = false
	grpStatus.duartion = 0
	grpStatus.voteSkip = []int64{}
}

func finalizeVote(chatID, msgID int64, ptcpCount int32) {
	percentage := float64(len(grpStatus.voteSkip)) / float64(ptcpCount) * 100

	status := "Failed"
	if percentage >= config.GetSuccessRate() {
		status = "Succeed"
	}

	format := fmt.Sprintf("Skip a song?\n"+
		"Vote count: %v\n"+
		"Vote Ended!\n\n"+
		"Status: %v", len(grpStatus.voteSkip), status)
	msg := helper.NewEditMessageText(chatID, msgID, nil, helper.NewFormattedText(format, nil))
	_, _ = bot.EditMessageText(msg)

	resetVote()
	if status == "Succeed" {
		fb2k.SetKillSwitch()
		if len(GetQueue()) == 0 {
			fb2k.PlayNext()
		} else {
			fb2k.PlaySelected(GetQueue()[0])
		}
	}
}

func endVote(chatID, msgID int64) {
	vs := grpStatus
	format := fmt.Sprintf("Skip a song?\n"+
		"Vote count: %v\n"+
		"Vote Ended!\n\n"+
		"Status: Generating vote results...", len(vs.voteSkip))
	msg := helper.NewEditMessageText(chatID, msgID, nil, helper.NewFormattedText(format, nil))
	_, _ = bot.EditMessageText(msg)

	if !config.IsWebEnabled() {
		c, err := userBot.GetChat(&tdlib.GetChatRequest{ChatId: chatID})
		if err != nil {
			resetVote()
			log.Println(err)
			return
		}
		if c.VideoChat.GroupCallId == 0 {
			resetVote()
			log.Println("No group call currently.")
			return
		}
		vc, err := userBot.GetGroupCall(&tdlib.GetGroupCallRequest{GroupCallId: c.VideoChat.GroupCallId})
		if err != nil {
			resetVote()
			log.Println(err)
			return
		}
		finalizeVote(chatID, msgID, vc.ParticipantCount)
	} else {
		finalizeVote(chatID, msgID, int32(len(grpStatus.Ptcps)))
	}
}

func setUserVote(chatID, msgID int64, userID int64, queryID tdlib.JsonInt64) {
	if config.IsJoinNeeded() {
		cm, err := bot.GetChatMember(&tdlib.GetChatMemberRequest{ChatId: config.GetChatID(), MemberId: helper.NewMessageSenderUser(userID)})
		if err != nil {
			bot.AnswerCallbackQuery(helper.NewAnswerCallbackQuery(queryID, "Failed to fetch chat info! Please try again later~", true, "", 10))
			return
		}

		if cm.Status.ChatMemberStatusType() == "chatMemberStatusLeft" {
			bot.AnswerCallbackQuery(helper.NewAnswerCallbackQuery(queryID, "Only users which are in the group can vote!", true, "", 10))
			return
		}
	}

	if utils.ContainsInt64(grpStatus.voteSkip, userID) {
		bot.AnswerCallbackQuery(helper.NewAnswerCallbackQuery(queryID, "You're already vote!", false, "", 45))
		return
	}

	hashedID := getUserIDHash(int64(userID))
	if !utils.ContainsString(GetPtcps(), hashedID) && config.IsPtcpsOnly() {
		bot.AnswerCallbackQuery(helper.NewAnswerCallbackQuery(queryID, "Only users which are in a voice chat can vote!", false, "", 5))
		return
	}

	AddVote(userID)
	updateVote(chatID, msgID, false)
}

// AddVote add user to vote list
func AddVote(userID int64) {
	if !utils.ContainsInt64(grpStatus.voteSkip, userID) {
		grpStatus.voteSkip = append(grpStatus.voteSkip, userID)
	}
}
