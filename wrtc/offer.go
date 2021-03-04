package wrtc

import (
	"strconv"
	"strings"

	"github.com/c0re100/go-tdlib"
	"github.com/pion/webrtc/v2"
)

type data struct {
	UFrag string
	Pwd   string
	Port  string
	Ssrc  int64
	Cert  []webrtc.Certificate
	Offer string
}

func extractDesc(pc *webrtc.PeerConnection, sdp string) *data {
	var lines []string
	if sdp != "" {
		lines = strings.Split(sdp, "\n")
	} else if peerConnection.LocalDescription() != nil {
		lines = strings.Split(peerConnection.LocalDescription().SDP, "\n")
	} else if peerConnection.PendingLocalDescription() != nil {
		lines = strings.Split(peerConnection.PendingLocalDescription().SDP, "\n")
	}
	var ufrag, pwd, port string
	var ssrc int64

	for _, s := range lines {
		if strings.Contains(s, "a=ice-ufrag:") {
			ufrag = strings.Split(s, "a=ice-ufrag:")[1]
		}
		if strings.Contains(s, "a=ice-pwd:") {
			pwd = strings.Split(s, "a=ice-pwd:")[1]
		}
		if strings.Contains(s, "a=ssrc:") {
			ssrc, _ = strconv.ParseInt(strings.Split(strings.Split(s, "a=ssrc:")[1], " ")[0], 10, 64)
		}
		if strings.Contains(s, "a=candidate:foundation 1 udp 2130706431 192.168.0.100 ") {
			port = strings.Split(strings.Split(s, "a=candidate:foundation 1 udp 2130706431 192.168.0.100 ")[1], " typ")[0]
		}
	}

	return &data{
		UFrag: ufrag,
		Pwd:   pwd,
		Port:  port,
		Ssrc:  ssrc,
		Cert:  pc.GetConfiguration().Certificates,
		Offer: pc.PendingLocalDescription().SDP,
	}
}

func createOfferSdp(resp *tdlib.GroupCallJoinResponse, ssrc string) string {
	var offerSdp string

	offerSdp += "v=0\n"
	offerSdp += "o=- 6543245 2 IN IP4 0.0.0.0\n"
	offerSdp += "s=-\n"
	offerSdp += "t=0 0\n"
	offerSdp += "a=group:BUNDLE 0\n"
	offerSdp += "a=ice-lite\n"
	offerSdp += "m=audio 1 UDP/TLS/RTP/SAVPF 111\n"
	offerSdp += "c=IN IP4 0.0.0.0\n"
	offerSdp += "a=mid:0\n"
	if resp != nil {
		offerSdp += "a=ice-ufrag:" + resp.Payload.Ufrag + "\n"
		offerSdp += "a=ice-pwd:" + resp.Payload.Pwd + "\n"
		for _, f := range resp.Payload.Fingerprints {
			offerSdp += "a=fingerprint:sha-256 " + f.Fingerprint + "\n"
		}
		offerSdp += "a=setup:passive\n"
		for i, c := range resp.Candidates {
			offerSdp += "a=candidate:" + strconv.Itoa(i+1) + " 1 udp " + c.Priority + " " + c.Ip + " " + c.Port + " typ host generation 0\n"
		}
	}
	offerSdp += "a=rtpmap:111 opus/48000/2\n"
	offerSdp += "a=rtcp-fb:111 transport-cc\n"
	offerSdp += "a=fmtp:111 minptime=10; useinbandfec=1\n"
	offerSdp += "a=rtcp:1 IN IP4 0.0.0.0\n"
	offerSdp += "a=rtcp-mux\n"
	offerSdp += "a=recvonly\n"

	return offerSdp
}

func createLocalSdp(resp *tdlib.GroupCallJoinResponse, offer string) string {
	var offerSdp string

	o := strings.Split(offer, "\n")
	for _, s := range o {
		if strings.Contains(s, "a=fingerprint:sha-256") {
			for _, f := range resp.Payload.Fingerprints {
				offerSdp += "a=fingerprint:sha-256 " + f.Fingerprint + "\n"
			}
		} else if strings.Contains(s, "a=ice-ufrag:") {
			offerSdp += "a=ice-ufrag:" + resp.Payload.Ufrag + "\n"
		} else if strings.Contains(s, "a=ice-pwd:") {
			offerSdp += "a=ice-pwd:" + resp.Payload.Pwd + "\n"
		} else {
			offerSdp += s + "\n"
		}
	}

	return offerSdp
}
