package data

import "encoding/json"

type Proxy struct {
	Name    string
	Address Address
	Kind    string
	Secret  string
}

func (p *Proxy) MarshalJSON() ([]byte, error) {
	aux := struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Kind    string `json:"kind"`
	}{
		Name:    p.Name,
		Address: p.Address.String(),
		Kind:    p.Kind,
	}
	return json.Marshal(aux)
}
