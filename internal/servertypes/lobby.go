package servertypes

import (
	"GateEngine/internal/data"
	"errors"
)

type lobby struct {
	lobbies  []data.Server
	TypeName TypeName
}

func (l lobby) GetServers() []data.Server {
	return l.lobbies
}

func (l lobby) GetTypeName() TypeName {
	return l.TypeName
}

func (l *lobby) AddIfNotExistAndSameType(server data.Server) error {
	if l.IsType(server) && !l.ServerExists(server) {
		return l.AddServer(server)
	}
	return errors.New("server already exists and/or is not of the same type")
}

func (l *lobby) AddServer(server data.Server) error {
	l.lobbies = append(l.lobbies, server)
	return nil
}

func (l lobby) IsType(server data.Server) bool {
	return server.Kind == string(l.TypeName)
}

func (l *lobby) RemoveServer(name string) error {
	var servers []data.Server
	var removedServers int
	for _, s := range l.lobbies {
		if s.Name != name {
			servers = append(servers, s)
		} else {
			removedServers++
		}
	}
	switch removedServers {
	case 0:
		return errors.New("server does not exist")
	case 1:
		l.lobbies = servers
		return nil
	default:
		return errors.New("multiple servers with the same name exist")
	}
}

func (l lobby) ServerExists(server data.Server) bool {
	for _, s := range l.lobbies {
		if s.Name == server.Name {
			return true
		}
	}
	return false
}
