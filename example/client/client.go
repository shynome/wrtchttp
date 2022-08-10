package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/shynome/wrtchttp/example/config"
	"github.com/shynome/wrtchttp/ortc"
	"github.com/shynome/wrtchttp/signaler/sdk"
)

func main() {
	pc := try.To1(
		ortc.New(nil, ortc.DefaultConfig))
	defer pc.Close()

	roffer := try.To1(
		exchangeOffer(ortc.Signal{}))

	try.To1(
		pc.HandleConnect(roffer))
	fmt.Println("wrtc connected")

	l := try.To1(
		net.Listen("tcp", ":7070"))

	for {
		conn := try.To1(l.Accept())
		go func(conn net.Conn) {
			if err := handleReq(pc, conn); err != nil {
				log.Println("handle req error: ", err)
				return
			}
		}(conn)
	}
}

var id uint16 = 0

func handleReq(pc *ortc.ORTC, conn net.Conn) (err error) {
	defer err2.Return(&err)

	id++

	dc := try.To1(
		pc.NewDataChannel(&id))
	defer dc.Close()

	rconn := try.To1(
		dc.Detach())

	go io.Copy(rconn, conn)
	io.Copy(conn, rconn)

	return
}

func exchangeOffer(offer ortc.Signal) (roffer ortc.Signal, err error) {
	defer err2.Return(&err)

	sdk := sdk.New(config.SignalerServer)

	b := try.To1(
		json.Marshal(offer))
	output := try.To1(
		sdk.Call(b))

	try.To(
		json.Unmarshal(output, &roffer))

	return
}
