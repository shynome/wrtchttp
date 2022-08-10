package ortc

import (
	"context"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
)

type Signal webrtc.SessionDescription

type ORTC struct {
	api *webrtc.API
	pc  *webrtc.PeerConnection
}

func NewAPI() (api *webrtc.API) {
	settingEngine := webrtc.SettingEngine{}
	settingEngine.DetachDataChannels()

	api = webrtc.NewAPI(
		webrtc.WithSettingEngine(settingEngine),
	)
	return
}

func New(api *webrtc.API, config webrtc.Configuration) (ortc *ORTC, err error) {
	defer err2.Return(&err)

	if api == nil {
		api = NewAPI()
	}

	pc := try.To1(
		api.NewPeerConnection(config))

	ortc = &ORTC{
		api: api,
		pc:  pc,
	}

	return
}

func (o *ORTC) HandleConnect(offer Signal) (roffer Signal, err error) {
	defer err2.Return(&err)

	var (
		pc = o.pc
	)

	try.To(
		pc.SetRemoteDescription(webrtc.SessionDescription(offer)))

	answer := try.To1(
		pc.CreateAnswer(nil))

	gatherComplete := webrtc.GatheringCompletePromise(pc)

	try.To(
		pc.SetLocalDescription(answer))

	<-gatherComplete

	roffer = Signal(*pc.LocalDescription())

	return
}

func (o *ORTC) CreateOffer() (sdp Signal, err error) {
	defer err2.Return(&err)

	var (
		pc = o.pc
	)

	try.To(
		o.makeOfferWithCandidates())

	offer := try.To1(
		pc.CreateOffer(nil))

	try.To(
		pc.SetLocalDescription(offer))
	sdp = Signal(offer)

	return
}

func (o *ORTC) makeOfferWithCandidates() (err error) {
	defer err2.Return(&err)

	var (
		pc = o.pc
	)

	dc := try.To1(
		pc.CreateDataChannel("_for_collect_candidates", nil))
	defer dc.Close()

	ctx, done := context.WithCancel(context.Background())
	pc.OnNegotiationNeeded(func() {
		done()
	})
	<-ctx.Done()

	return
}

func (o *ORTC) Handshake(offer Signal) (err error) {
	defer err2.Return(&err)
	var pc = o.pc

	try.To(
		pc.SetRemoteDescription(webrtc.SessionDescription(offer)))

	return
}

func (o *ORTC) NewDataChannel(id *uint16) (dc *webrtc.DataChannel, err error) {
	defer err2.Return(&err)

	var pc = o.pc

	params := &webrtc.DataChannelInit{
		ID: id,
	}

	dc = try.To1(
		pc.CreateDataChannel("", params))
	// wait dc opened
	w := make(chan struct{})
	dc.OnOpen(func() { close(w) })
	<-w

	return
}

func (o *ORTC) OnDataChannel(f func(dc *webrtc.DataChannel)) {
	o.pc.OnDataChannel(f)
}

func (o *ORTC) PC() *webrtc.PeerConnection {
	return o.pc
}

// type Stop interface{ Stop() error }
// func stop(stoper Stop) { try.To(stoper.Stop()) }

func (o *ORTC) Close() (err error) {
	defer err2.Return(&err)
	o.pc.Close()
	return
}
