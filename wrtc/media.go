package wrtc

import (
	"fmt"
	"math/rand"

	"github.com/pion/webrtc/v2"
)

func setupMedia() {
	id := rand.Uint32()
	for {
		if id >= 2147483647 {
			id = rand.Uint32()
		} else {
			break
		}
	}
	audioTrack, addTrackErr := peerConnection.NewTrack(getPayloadType(mediaEngine, webrtc.RTPCodecTypeAudio, "opus"), id, "audio", "radio")
	if addTrackErr != nil {
		panic(addTrackErr)
	}
	if _, addTrackErr = peerConnection.AddTransceiverFromTrack(audioTrack, webrtc.RtpTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionSendrecv,
	}); addTrackErr != nil {
		panic(addTrackErr)
	}
}

func getPayloadType(m webrtc.MediaEngine, codecType webrtc.RTPCodecType, codecName string) uint8 {
	for _, codec := range m.GetCodecsByKind(codecType) {
		if codec.Name == codecName {
			return codec.PayloadType
		}
	}
	panic(fmt.Sprintf("Remote peer does not support %s", codecName))
}
