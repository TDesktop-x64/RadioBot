package wrtc

import (
	"fmt"
	"strconv"

	"github.com/c0re100/RadioBot/config"
	// "github.com/c0re100/go-tdlib"
	"github.com/pion/webrtc/v2"
)

// Connect connect to group call server
func Connect(resp string, d *Data) {
	rSdp := createOfferSdp(resp, strconv.FormatInt(d.Ssrc, 10))
	rOffer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  rSdp,
	}
	err := peerConnection.SetRemoteDescription(rOffer)
	if err != nil {
		panic(err)
	}

	_, err = peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	select {
	case <-closeRTC:
		peerConnection.Close()
		fmt.Println("WebRTC connection closed.")
	default:

	}
}

// Disconnect disconnect all connections
func Disconnect() {
	if !config.IsWebEnabled() {
		closeRTC <- true
		c, _ := userBot.GetChat(config.GetChatID())
		gc, _ := userBot.GetGroupCall(c.VoiceChat.GroupCallId)
		userBot.LeaveGroupCall(gc.Id)
	}
}

// GetConnection get peer connection
func GetConnection() *webrtc.PeerConnection {
	return peerConnection
}

// GetCurrentSDP get SDP string
func GetCurrentSDP() string {
	if peerConnection.LocalDescription() != nil {
		return peerConnection.LocalDescription().SDP
	} else if peerConnection.PendingLocalDescription() != nil {
		return peerConnection.PendingLocalDescription().SDP
	} else {
		return ""
	}
}
