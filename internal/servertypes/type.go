package servertypes

import "GateEngine/internal/data"

type TypeName string

type ServerType interface {
	GetServers() []data.Server
	GetTypeName() TypeName
	AddIfNotExistAndSameType(server data.Server) error
	AddServer(server data.Server) error
	IsType(server data.Server) bool
	ServerExists(server data.Server) bool
	RemoveServer(name string) error
	[]data.Server
	TypeName
}

type LobbyServerType interface {
}
