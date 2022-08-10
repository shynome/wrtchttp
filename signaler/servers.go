package signaler

import (
	"errors"
	"sync"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type Servers struct {
	mu      *sync.RWMutex
	servers map[string]*Server
}

func NewServers() *Servers {
	return &Servers{
		mu:      &sync.RWMutex{},
		servers: map[string]*Server{},
	}
}

type Server struct {
	Name     string
	Password string
}

func (s *Servers) Get(auth Auth) (server *Server, err error) {
	defer err2.Return(&err)

	try.To(
		checkAuth(auth))

	var (
		name = auth.Server
		pass = auth.Password
	)
	s.mu.RLock()
	server, ok := s.servers[name]
	s.mu.RUnlock()
	if !ok || server == nil {
		server = s.Create(name, pass)
		return
	}
	if server.Password != pass {
		err = errors.New("server password is incorrect")
		return
	}
	return
}

func (s *Servers) Create(name, pass string) (server *Server) {
	s.mu.Lock()
	defer s.mu.Unlock()
	server = &Server{
		Name:     name,
		Password: pass,
	}
	s.servers[name] = server
	return
}

func checkAuth(auth Auth) error {
	var (
		name = auth.Server
		pass = auth.Password
	)
	if len(name) > 64 || len(pass) > 64 {
		return errors.New("server name or password is too long")
	}
	if name == "" {
		return errors.New("server name is empty")
	}
	if pass == "" {
		return errors.New("server password is empty")
	}
	return nil
}
