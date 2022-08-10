package ortc

import "github.com/pion/webrtc/v3"

var DefaultConfig = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{URLs: []string{"stun:stun.l.google.com:19302"}},
	},
}
