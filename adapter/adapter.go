package adapter

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/donovanhide/eventsource"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/pion/datachannel"
	"github.com/pion/webrtc/v3"
	"github.com/shynome/wrtchttp/ortc"
	"github.com/shynome/wrtchttp/signaler/sdk"
)

type Adapter interface {
	Accept() (conn datachannel.ReadWriteCloser, err error)
}

type DefaultAdapter struct {
	wrtcApi *webrtc.API
	ice     webrtc.ICEGatherOptions

	sdk *sdk.Sdk
	// key string

	conns chan datachannel.ReadWriteCloser
}

var _ Adapter = &DefaultAdapter{}

func NewAdapter(signaler string) (adapter *DefaultAdapter, err error) {
	defer err2.Return(&err)
	api := ortc.NewAPI()
	sdk := sdk.New(signaler)
	adapter = &DefaultAdapter{
		wrtcApi: api,
		ice:     ortc.DefaultICEGatherOptions,
		sdk:     sdk,
		conns:   make(chan datachannel.ReadWriteCloser),
	}
	return
}

func (a *DefaultAdapter) Accept() (conn datachannel.ReadWriteCloser, err error) {
	conn, ok := <-a.conns
	if !ok {
		return nil, errors.New("Closed")
	}
	return conn, nil
}

func (a *DefaultAdapter) Listen() {
	stream := try.To1(eventsource.Subscribe(a.sdk.Endpoint, ""))
	for ev := range stream.Events {
		go a.handleIncomingRequest(ev)
	}
}

func (a *DefaultAdapter) handleIncomingRequest(ev eventsource.Event) (err error) {
	defer err2.Return(&err)

	pc := try.To1(
		ortc.New(a.wrtcApi, a.ice))
	a.onDataChannel(pc)

	var offer ortc.Signal
	try.To(
		json.Unmarshal([]byte(ev.Data()), &offer))

	roffer := try.To1(
		json.Marshal(pc.Signal))
	try.To(
		a.sdk.Dial(ev.Id(), roffer))

	try.To(
		pc.HandShake(offer, webrtc.ICERoleControlled))

	return
}

func (a *DefaultAdapter) onDataChannel(pc *ortc.ORTC) {
	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		dc.OnOpen(func() {
			rw, err := dc.Detach()
			if err != nil {
				log.Println("got datachannel detach failed")
				return
			}
			a.conns <- rw
		})
	})
}
