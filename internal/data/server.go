package data

import (
	"encoding/json"
	"fmt"
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
	"time"
)

type status struct {
	Description chat.Message
	Players     struct {
		Max    int
		Online int
		Sample []struct {
			ID   uuid.UUID
			Name string
		}
	}
	Version struct {
		Name     string
		Protocol int
	}
	Favicon Icon
	Delay   time.Duration
}

type Icon string

type Address struct {
	IP   string
	Port int
}

type Server struct {
	Name        string  `json:"name"`
	Address     Address `json:"address"`
	Fallback    bool    `json:"fallback"`
	Kind        string  `json:",omitzero"`
	MaxPlayers  int     `json:",omitzero"`
	PlayerCount int     `json:",omitzero"`
}

func (s *Server) UpdatePlayerCount() error {
	response, _, err := bot.PingAndListTimeout(fmt.Sprintf("%s:%d", s.Address.IP, s.Address.Port), time.Second*5)
	if err != nil {
		return err
	}
	var status status
	err = json.Unmarshal(response, &status)
	if err != nil {
		return err
	}
	s.MaxPlayers = status.Players.Max
	s.PlayerCount = status.Players.Online
	return nil
}

func (s *Server) MarshalJSON() ([]byte, error) {
	aux := struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Fallback bool   `json:"fallback"`
	}{
		Name:     s.Name,
		Address:  fmt.Sprintf("%s:%d", s.Address.IP, s.Address.Port),
		Fallback: s.Fallback,
	}
	return json.Marshal(aux)
}
