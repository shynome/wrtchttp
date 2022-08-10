module github.com/shynome/wrtchttp

go 1.18

require (
	github.com/donovanhide/eventsource v0.0.0-20210830082556-c59027999da0
	github.com/gorilla/websocket v1.5.0
	github.com/lainio/err2 v0.8.7
	github.com/pion/datachannel v1.5.2
	github.com/pion/webrtc/v3 v3.1.43
	github.com/rs/xid v1.4.0
)

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/pion/dtls/v2 v2.1.5 // indirect
	github.com/pion/ice/v2 v2.2.6 // indirect
	github.com/pion/interceptor v0.1.11 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/mdns v0.0.5 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtcp v1.2.9 // indirect
	github.com/pion/rtp v1.7.13 // indirect
	github.com/pion/sctp v1.8.2 // indirect
	github.com/pion/sdp/v3 v3.0.5 // indirect
	github.com/pion/srtp/v2 v2.0.10 // indirect
	github.com/pion/stun v0.3.5 // indirect
	github.com/pion/transport v0.13.1 // indirect
	github.com/pion/turn/v2 v2.0.8 // indirect
	github.com/pion/udp v0.1.1 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/net v0.0.0-20220802222814-0bcc04d9c69b // indirect
	golang.org/x/sys v0.0.0-20220731174439-a90be440212d // indirect
)

replace github.com/pion/sctp v1.8.2 => github.com/shynome/pion-sctp v0.0.0-20220809083843-2d8ed9b66aa1
