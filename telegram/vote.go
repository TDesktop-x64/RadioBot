package telegram

import (
	"fmt"
	"log"
	"time"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/fb2k"
	"github.com/c0re100/RadioBot/utils"
	"github.com/c0re100/go-tdlib"
)

type groupStatus struct {
	chatID       int64
	msgID        int64
	vcID         int32
	duartion     int32
	Ptcps        []int32
	voteSkip     []int32
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

func startVote(chatID, msgID int64, userID int32) {
	if chatID != config.GetChatID() {
		return
	}

	if !config.IsVoteEnabled() {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("This group is not allowed to vote.", nil), true, true)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		return
	}

	if !config.IsWebEnabled() {
		c, err := userBot.GetChat(chatID)
		if err != nil {
			log.Println(err)
			return
		}

		if c.VoiceChatGroupCallId == 0 {
			msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("This group do not have a voice chat.", nil), true, true)
			bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
			return
		}
		// Preload all users
		_, _ = userBot.LoadGroupCallParticipants(c.VoiceChatGroupCallId, 5000)
	}

	if !utils.ContainsInt32(grpStatus.Ptcps, userID) {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Only users which are in a voice chat can vote!", nil), true, true)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		return
	}

	if grpStatus.isVoting {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Vote in progress...", nil), true, true)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		return
	}

	if time.Now().Unix() < grpStatus.lastVoteTime+config.GetReleaseTime() {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Skip a song was voted too recently...", nil), true, true)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		return
	}

	voteKb := tdlib.NewReplyMarkupInlineKeyboard([][]tdlib.InlineKeyboardButton{
		{
			*tdlib.NewInlineKeyboardButton("Yes - 1", tdlib.NewInlineKeyboardButtonTypeCallback([]byte("vote_skip"))),
		},
	})

	msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Skip a song?", nil), true, true)
	m, err := bot.SendMessage(chatID, 0, msgID, nil, voteKb, msgText)
	if err != nil {
		log.Println("Can't send message.")
		return
	}
	grpStatus.isVoting = true
	grpStatus.duartion = config.GetVoteTime()
	grpStatus.msgID = m.Id
	grpStatus.lastVoteTime = time.Now().Unix()

	if !utils.ContainsInt32(grpStatus.voteSkip, userID) {
		grpStatus.voteSkip = append(grpStatus.voteSkip, userID)
	}
	updateVote(chatID, m.Id, false)
	addVoteJob(chatID, m.Id)
	// Wait 15 seconds
	time.Sleep(15 * time.Second)
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
	voteKb := tdlib.NewReplyMarkupInlineKeyboard([][]tdlib.InlineKeyboardButton{
		{
			*tdlib.NewInlineKeyboardButton(fmt.Sprintf("Yes - %v", len(grpStatus.voteSkip)), tdlib.NewInlineKeyboardButtonTypeCallback([]byte("vote_skip"))),
		},
	})

	msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText(fmt.Sprintf("Skip a song?\n"+
		"Vote count: %v\n"+
		"Vote timeleft: %v second(s)", len(grpStatus.voteSkip), grpStatus.duartion), nil), true, true)
	bot.EditMessageText(chatID, msgID, voteKb, msgText)
}

func resetVote() {
	grpStatus.isLoadPtcps = false
	grpStatus.isVoting = false
	grpStatus.duartion = 0
	grpStatus.voteSkip = []int32{}
}

func finalizeVote(chatID, msgID int64, ptcpCount int32) {
	percentage := float64(len(grpStatus.voteSkip)) / float64(ptcpCount) * 100

	status := "Failed"
	if percentage >= config.GetSuccessRate() {
		status = "Succeed"
	}

	msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText(fmt.Sprintf("Skip a song?\n"+
		"Vote count: %v\n"+
		"Vote Ended!\n\n"+
		"Status: %v", len(grpStatus.voteSkip), status), nil), true, true)
	bot.EditMessageText(chatID, msgID, nil, msgText)

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
	msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText(fmt.Sprintf("Skip a song?\n"+
		"Vote count: %v\n"+
		"Vote Ended!\n\n"+
		"Status: Generating vote results...", len(vs.voteSkip)), nil), true, true)
	bot.EditMessageText(chatID, vs.msgID, nil, msgText)

	if !config.IsWebEnabled() {
		c, err := userBot.GetChat(chatID)
		if err != nil {
			resetVote()
			log.Println(err)
			return
		}
		if c.VoiceChatGroupCallId == 0 {
			resetVote()
			log.Println("No group call currently.")
			return
		}
		vc, err := userBot.GetGroupCall(c.VoiceChatGroupCallId)
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

func setUserVote(chatID, msgID int64, userID int32, queryID tdlib.JSONInt64) {
	if config.IsJoinNeeded() {
		cm, err := bot.GetChatMember(config.GetChatID(), userID)
		if err != nil {
			bot.AnswerCallbackQuery(queryID, "Failed to fetch chat info! Please try again later~", true, "", 10)
			return
		}

		if cm.Status.GetChatMemberStatusEnum() == "chatMemberStatusLeft" {
			bot.AnswerCallbackQuery(queryID, "Only users which are in the group can vote!", true, "", 10)
			return
		}
	}

	if utils.ContainsInt32(grpStatus.voteSkip, userID) {
		bot.AnswerCallbackQuery(queryID, "You're already vote!", false, "", 45)
		return
	}

	if config.IsPtcpsOnly() {
		bot.AnswerCallbackQuery(queryID, "Only users which are in a voice chat can vote!", false, "", 5)
		return
	}

	AddVote(userID)
	updateVote(chatID, msgID, false)
}

// AddVote add user to vote list
func AddVote(userID int32) {
	if !utils.ContainsInt32(grpStatus.voteSkip, userID) {
		grpStatus.voteSkip = append(grpStatus.voteSkip, userID)
	}
}
