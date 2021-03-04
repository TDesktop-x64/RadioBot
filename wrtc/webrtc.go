package wrtc

import (
	"fmt"
	"log"

	"github.com/c0re100/go-tdlib"
	"github.com/pion/webrtc/v2"
)

var (
	peerConnection *webrtc.PeerConnection
	mediaEngine    webrtc.MediaEngine
	closeRTC       = make(chan bool, 1)
	userBot        *tdlib.Client
)

func setup() {
	mediaEngine = webrtc.MediaEngine{}
	peerConnection = &webrtc.PeerConnection{}
	mediaEngine.RegisterCodec(webrtc.NewRTPOpusCodec(111, 48000))
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))
	peerConnection, _ = api.NewPeerConnection(webrtc.Configuration{})
	setupMedia()
}

func CreateOffer(bot *tdlib.Client) *data {
	setup()
	userBot = bot

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("Connection State has changed %s \n", connectionState.String())
	})

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i != nil {
			fmt.Println(i.ToJSON())
		}
	})

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	return extractDesc(peerConnection, "")
}
