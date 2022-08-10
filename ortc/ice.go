package ortc

import "github.com/pion/webrtc/v3"

var DefaultICEGatherOptions = webrtc.ICEGatherOptions{
	ICEServers: []webrtc.ICEServer{
		{URLs: []string{"stun:stun.l.google.com:19302"}},
	},
}
