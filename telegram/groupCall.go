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
	gc, _ := userBot.GetGroupCall(c.VoiceChatGroupCallId)
	grpStatus.vcID = gc.Id

	data := wrtc.CreateOffer(userBot)
	payload := tdlib.NewGroupCallPayload(data.UFrag, data.Pwd, nil)
	for _, c := range data.Cert {
		fp, _ := c.GetFingerprints()
		for _, f := range fp {
			payload.Fingerprints = append(payload.Fingerprints, *tdlib.NewGroupCallPayloadFingerprint(f.Value, "active", f.Value))
		}
	}
	gcResp, err := userBot.JoinGroupCall(gc.Id, payload, int32(data.Ssrc), false)
	if err != nil {
		log.Println(err)
		return
	}

	addLoadGroupCallPtpcsJob()
	if !sch.IsRunning() {
		log.Println("Starting scheduler...")
		startScheduler()
	}

	go wrtc.Connect(gcResp, data)
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

func newGroupCallPtcpUpdate() {
	fmt.Println("[Music] New GroupCallParticipant Receiver")
	eventFilter := func(msg *tdlib.TdMessage) bool {
		return true
	}

	receiver := userBot.AddEventReceiver(&tdlib.UpdateGroupCallParticipant{}, eventFilter, 5000)
	for newMsg := range receiver.Chan {
		updateMsg := (newMsg).(*tdlib.UpdateGroupCallParticipant)
		gcID := updateMsg.GroupCallId
		userID := updateMsg.Participant.UserId
		if grpStatus.vcID == gcID {
			if updateMsg.Participant.Order == 0 {
				if userID == userBotID && wrtc.GetConnection().ConnectionState().String() != "closed" {
					time.Sleep(1 * time.Second)
					log.Println("Userbot left voice chat...re-join now!")
					joinGroupCall()
				}
				RemovePtcp(userID)
			}
			AddPtcp(userID)
		}
	}
}

// AddPtcp add user to participant list
func AddPtcp(userID int32) {
	if !utils.ContainsInt32(grpStatus.Ptcps, userID) {
		//log.Printf("User %v joined voice chat.\n", uId)
		grpStatus.Ptcps = append(grpStatus.Ptcps, userID)
	}
}

// RemovePtcp remove user from participant list
func RemovePtcp(userID int32) {
	//log.Printf("User %v left voice chat.\n", uId)
	grpStatus.Ptcps = utils.FilterInt32(grpStatus.Ptcps, func(s int32) bool {
		return s != userID
	})
}

// ResetPtcps reset participant list
func ResetPtcps() {
	grpStatus.Ptcps = []int32{}
}

// GetPtcps get participant list
func GetPtcps() []int32 {
	return grpStatus.Ptcps
}
