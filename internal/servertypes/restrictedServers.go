package servertypes

import "GateEngine/internal/data"

type RestrictedServers interface {
	IsAllowed(player data.Player, server data.Server) (bool, error)
}
