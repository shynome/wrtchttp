package signaler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/donovanhide/eventsource"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/rs/xid"
	"github.com/shynome/wrtchttp/events"
	"github.com/shynome/wrtchttp/signaler/sdk"
)

type Signaler struct {
	ess     *eventsource.Server
	evs     *events.Events[[]byte]
	servers *Servers
}

var _ http.Handler = &Signaler{}

func New() (s *Signaler) {
	s = &Signaler{
		ess:     eventsource.NewServer(),
		evs:     events.New[[]byte](),
		servers: NewServers(),
	}
	return
}

func (s *Signaler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpf(s.selectDialer)(w, r)
}

func (s *Signaler) selectDialer(w http.ResponseWriter, r *http.Request) (err error) {
	switch r.Method {
	case http.MethodGet:
		return s.subscribe(w, r)
	case http.MethodPost:
		return s.handleCall(w, r)
	case http.MethodDelete:
		return s.handleDial(w, r)
	default:
		http.Error(w, fmt.Sprintf("deny method: %s \r\n", r.Method), 400)
	}
	return
}

func (s *Signaler) subscribe(w http.ResponseWriter, r *http.Request) (err error) {
	defer err2.Return(&err)
	auth := getAuth(r.URL)
	server := try.To1(
		s.servers.Get(auth))
	s.ess.Handler(server.Name)(w, r)
	return
}

func (s *Signaler) handleCall(w http.ResponseWriter, r *http.Request) (err error) {
	defer err2.Return(&err)
	auth := getAuth(r.URL)
	server := try.To1(
		s.servers.Get(auth))

	id := xid.New().String()
	result := s.evs.On(id)
	defer s.evs.Off(id)

	b := try.To1(
		io.ReadAll(r.Body))
	var ev = &sdk.Event{ID: id, Body: b}

	s.ess.Publish([]string{server.Name}, ev)

	rbody := <-result

	h := w.Header()
	h.Set("Content-Type", "application/octet-stream")
	try.To1(
		io.Copy(w, bytes.NewReader(rbody)))

	return
}

func (s *Signaler) handleDial(w http.ResponseWriter, r *http.Request) (err error) {
	defer err2.Return(&err)
	id := r.Header.Get("X-Event-Id")
	b := try.To1(
		io.ReadAll(r.Body))

	fmt.Println("result ", id)
	try.To(
		s.evs.Emit(id, b))
	fmt.Println("finish ", id)

	return
}

type Auth struct {
	Server   string
	Password string
}

func getAuth(u *url.URL) Auth {
	q := u.Query()
	return Auth{q.Get("server"), q.Get("password")}
}
