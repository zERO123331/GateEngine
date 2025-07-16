package servertypes

import "GateEngine/internal/data"

type minigamesLobby struct {
}

func (m minigamesLobby) GetServers() []data.Server {
	// TODO implement me
	panic("implement me")
}

func (m minigamesLobby) GetTypeName() TypeName {
	// TODO implement me
	panic("implement me")
}

func (m minigamesLobby) AddIfNotExistAndSameType(server data.Server) error {
	// TODO implement me
	panic("implement me")
}

func (m minigamesLobby) AddServer(server data.Server) error {
	// TODO implement me
	panic("implement me")
}

func (m minigamesLobby) IsType(server data.Server) bool {
	// TODO implement me
	panic("implement me")
}

func (m minigamesLobby) ServerExists(server data.Server) bool {
	// TODO implement me
	panic("implement me")
}

func (m minigamesLobby) RemoveServer(name string) error {
	// TODO implement me
	panic("implement me")
}
