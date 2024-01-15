package telegram

import (
	"fmt"
	"log"
	"time"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/utils"
	"github.com/c0re100/RadioBot/wrtc"
	tdlib "github.com/c0re100/gotdlib/client"
)

func joinGroupCall() {
	c, _ := userBot.GetChat(&tdlib.GetChatRequest{ChatId: config.GetChatID()})
	gc, _ := userBot.GetGroupCall(&tdlib.GetGroupCallRequest{GroupCallId: c.VideoChat.GroupCallId})
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

func loadParticipants(chatID int64, userID int64) {
	if isAdmin(chatID, userID) {
		gc, _ := userBot.GetGroupCall(&tdlib.GetGroupCallRequest{GroupCallId: grpStatus.vcID})
		if gc.LoadedAllParticipants {
			return
		}
		_, _ = userBot.LoadGroupCallParticipants(&tdlib.LoadGroupCallParticipantsRequest{GroupCallId: gc.Id, Limit: 5000})
	}
}

func newGroupCallUpdate() {
	fmt.Println("[Music] New GroupCall Receiver")

	listener := userBot.AddEventReceiver(&tdlib.UpdateGroupCall{}, 100)
	for newMsg := range listener.Updates {
		updateMsg := (newMsg).(*tdlib.UpdateGroupCall)
		gcID := updateMsg.GroupCall.Id
		// todo
		if grpStatus.vcID == gcID && grpStatus.isLoadPtcps && updateMsg.GroupCall.LoadedAllParticipants {
			finalizeVote(grpStatus.chatID, grpStatus.msgID, updateMsg.GroupCall.ParticipantCount)
		}
	}
}

func GetsenderId(sender tdlib.MessageSender) int64 {
	if sender.MessageSenderType() == "messageSenderUser" {
		return sender.(*tdlib.MessageSenderUser).UserId
	} else {
		return sender.(*tdlib.MessageSenderChat).ChatId
	}
}

func newGroupCallPtcpUpdate() {
	fmt.Println("[Music] New GroupCallParticipant Receiver")

	listener := userBot.AddEventReceiver(&tdlib.UpdateGroupCallParticipant{}, 5000)
	for newMsg := range listener.Updates {
		updateMsg := (newMsg).(*tdlib.UpdateGroupCallParticipant)
		gcID := updateMsg.GroupCallId
		userID := GetsenderId(updateMsg.Participant.ParticipantId)
		if grpStatus.vcID == gcID {
			hashedID := getUserIDHash(userID)
			if updateMsg.Participant.Order == "0" {
				if userID == userBotID && wrtc.GetConnection().ConnectionState().String() != "closed" {
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
