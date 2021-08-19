package telegram

import (
	"fmt"
	"log"
	"time"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/utils"
	"github.com/c0re100/RadioBot/wrtc"
	"github.com/c0re100/go-tdlib"
)

func joinGroupCall() {
	c, _ := userBot.GetChat(config.GetChatID())
	gc, _ := userBot.GetGroupCall(c.VoiceChat.GroupCallId)
	grpStatus.vcID = gc.Id

	data := wrtc.CreateOffer(userBot)
	// todo
	//payload := tdlib.NewGroupCallPayload(data.UFrag, data.Pwd, nil)
	//for _, c := range data.Cert {
	//	fp, _ := c.GetFingerprints()
	//	for _, f := range fp {
	//		payload.Fingerprints = append(payload.Fingerprints, *tdlib.NewGroupCallPayloadFingerprint(f.Value, "active", f.Value))
	//	}
	//}
	//gcResp, err := userBot.JoinGroupCall(gc.Id, payload, int32(data.Ssrc), false)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	addLoadGroupCallPtpcsJob()
	if !sch.IsRunning() {
		log.Println("Starting scheduler...")
		startScheduler()
	}

	go wrtc.Connect("todo", data)
}

func loadParticipants(chatID int64, userID int32) {
	if isAdmin(chatID, userID) {
		gc, _ := userBot.GetGroupCall(grpStatus.vcID)
		if gc.LoadedAllParticipants {
			return
		}
		userBot.LoadGroupCallParticipants(gc.Id, 5000)
	}
}

func newGroupCallUpdate() {
	fmt.Println("[Music] New GroupCall Receiver")
	eventFilter := func(msg *tdlib.TdMessage) bool {
		return true
	}

	receiver := userBot.AddEventReceiver(&tdlib.UpdateGroupCall{}, eventFilter, 100)
	for newMsg := range receiver.Chan {
		updateMsg := (newMsg).(*tdlib.UpdateGroupCall)
		gcID := updateMsg.GroupCall.Id
		// todo
		if grpStatus.vcID == gcID && grpStatus.isLoadPtcps && updateMsg.GroupCall.LoadedAllParticipants {
			finalizeVote(grpStatus.chatID, grpStatus.msgID, updateMsg.GroupCall.ParticipantCount)
		}
	}
}

func GetSenderId(sender tdlib.MessageSender) int64 {
	if sender.GetMessageSenderEnum() == "messageSenderUser" {
		return int64(sender.(*tdlib.MessageSenderUser).UserId)
	} else {
		return sender.(*tdlib.MessageSenderChat).ChatId
	}
}

func newGroupCallPtcpUpdate() {
	fmt.Println("[Music] New GroupCallParticipant Receiver")
	eventFilter := func(msg *tdlib.TdMessage) bool {
		return true
	}

	receiver := userBot.AddEventReceiver(&tdlib.UpdateGroupCallParticipant{}, eventFilter, 5000)
	for newMsg := range receiver.Chan {
		updateMsg := (newMsg).(*tdlib.UpdateGroupCallParticipant)
		gcID := updateMsg.GroupCallId
		userID := GetSenderId(updateMsg.Participant.ParticipantId)
		if grpStatus.vcID == gcID {
			hashedID := getUserIDHash(userID)
			if updateMsg.Participant.Order == "0" {
				if userID == int64(userBotID) && wrtc.GetConnection().ConnectionState().String() != "closed" {
					time.Sleep(1 * time.Second)
					log.Println("Userbot left voice chat...re-join now!")
					joinGroupCall()
				}
				RemovePtcp(hashedID)
			}
			AddPtcp(hashedID)
		}
	}
}

// AddPtcp add user to participant list
func AddPtcp(hashedID string) {
	if !utils.ContainsString(grpStatus.Ptcps, hashedID) {
		//log.Printf("User %v joined voice chat.\n", uId)
		grpStatus.Ptcps = append(grpStatus.Ptcps, hashedID)
	}
}

// RemovePtcp remove user from participant list
func RemovePtcp(hashedID string) {
	//log.Printf("User %v left voice chat.\n", uId)
	grpStatus.Ptcps = utils.FilterString(grpStatus.Ptcps, func(s string) bool {
		return s != hashedID
	})
}

// ResetPtcps reset participant list
func ResetPtcps() {
	grpStatus.Ptcps = []string{}
}

// GetPtcps get participant list
func GetPtcps() []string {
	return grpStatus.Ptcps
}
