package ortc

import (
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
)

type Signal struct {
	ICECandidates    []webrtc.ICECandidate   `json:"iceCandidates"`
	ICEParameters    webrtc.ICEParameters    `json:"iceParameters"`
	DTLSParameters   webrtc.DTLSParameters   `json:"dtlsParameters"`
	SCTPCapabilities webrtc.SCTPCapabilities `json:"sctpCapabilities"`
}

type ORTC struct {
	Signal Signal

	api  *webrtc.API
	ice  *webrtc.ICETransport
	dtls *webrtc.DTLSTransport
	sctp *webrtc.SCTPTransport
}

func NewAPI() (api *webrtc.API) {
	settingEngine := webrtc.SettingEngine{}
	settingEngine.DetachDataChannels()

	api = webrtc.NewAPI(
		webrtc.WithSettingEngine(settingEngine),
	)
	return
}

func New(api *webrtc.API, iceConfig webrtc.ICEGatherOptions) (ortc *ORTC, err error) {
	defer err2.Return(&err)

	if api == nil {
		api = NewAPI()
	}

	gatherer := try.To1(
		api.NewICEGatherer(iceConfig))

	ice := api.NewICETransport(gatherer)

	dtls := try.To1(
		api.NewDTLSTransport(ice, nil))

	sctp := api.NewSCTPTransport(dtls)

	gatherFinished := make(chan struct{})
	gatherer.OnLocalCandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			close(gatherFinished)
		}
	})
	try.To(
		gatherer.Gather())
	<-gatherFinished

	iceCandidates := try.To1(
		gatherer.GetLocalCandidates())
	iceParams := try.To1(
		gatherer.GetLocalParameters())
	dtlsParams := try.To1(
		dtls.GetLocalParameters())
	sctpCapabilities := sctp.GetCapabilities()

	ortc = &ORTC{
		api:  api,
		ice:  ice,
		dtls: dtls,
		sctp: sctp,
	}

	ortc.Signal = Signal{
		ICECandidates:    iceCandidates,
		ICEParameters:    iceParams,
		DTLSParameters:   dtlsParams,
		SCTPCapabilities: sctpCapabilities,
	}

	return
}

func (o *ORTC) HandShake(signal Signal, iceRole webrtc.ICERole) (err error) {
	defer err2.Return(&err)

	var (
		ice  = o.ice
		dtls = o.dtls
		sctp = o.sctp
	)

	try.To(
		ice.SetRemoteCandidates(signal.ICECandidates))

	try.To(
		ice.Start(nil, signal.ICEParameters, &iceRole))

	try.To(
		dtls.Start(signal.DTLSParameters))

	try.To(
		sctp.Start(signal.SCTPCapabilities))

	return
}

func (o *ORTC) NewDataChannel(id *uint16) (dc *webrtc.DataChannel, err error) {
	defer err2.Return(&err)

	var api = o.api

	params := &webrtc.DataChannelParameters{
		ID: id,
	}

	dc = try.To1(
		api.NewDataChannel(o.sctp, params))

	return
}

func (o *ORTC) OnDataChannel(f func(dc *webrtc.DataChannel)) {
	o.sctp.OnDataChannel(f)
}

type Stop interface{ Stop() error }

func stop(stoper Stop) { try.To(stoper.Stop()) }

func (o *ORTC) Close() (err error) {
	defer err2.Return(&err)
	defer stop(o.ice)
	defer stop(o.dtls)
	defer stop(o.sctp)
	return
}
