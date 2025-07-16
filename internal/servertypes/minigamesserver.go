package servertypes

import (
	"GateEngine/internal/data"
	"errors"
)

type MinigamesServerType struct {
	servers  []data.Server
	TypeName TypeName
}

func (m MinigamesServerType) IsAllowed(player data.Player, server data.Server) (bool, error) {
	if m.IsType(server) && m.ServerExists(server) {
		return true, nil
	}
	return false, errors.New("server is not of the same type or registered")
}

func (m MinigamesServerType) GetServers() []data.Server {
	return m.servers
}

func (m MinigamesServerType) GetTypeName() TypeName {
	return m.TypeName
}

func (m *MinigamesServerType) AddIfNotExistAndSameType(server data.Server) error {
	if m.IsType(server) && !m.ServerExists(server) {
		return m.AddServer(server)
	}
	return errors.New("server already exists and/or is not of the same type")
}

func (m *MinigamesServerType) AddServer(server data.Server) error {
	m.servers = append(m.servers, server)
	return nil
}

func (m MinigamesServerType) IsType(server data.Server) bool {
	return server.Kind == string(m.TypeName)
}

func (m MinigamesServerType) ServerExists(server data.Server) bool {
	for _, s := range m.servers {
		if s.Name == server.Name {
			return true
		}
	}
	return false
}

func (m *MinigamesServerType) RemoveServer(name string) error {
	var servers []data.Server
	var removedServers int
	for _, s := range m.servers {
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
		m.servers = servers
		return nil
	default:
		return errors.New("multiple servers with the same name exist")
	}
}
