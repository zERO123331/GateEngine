package data

import "encoding/json"

type Proxy struct {
	Name       string
	Address    Address
	Kind       string
	ID         int
	APIAddress Address
	Secret     string
}

func (p *Proxy) MarshalJSON() ([]byte, error) {
	aux := struct {
		Name       string `json:"name"`
		Address    string `json:"address"`
		Kind       string `json:"kind"`
		ID         int    `json:"id"`
		APIAddress string `json:"apiAddress"`
		Secret     string `json:"secret"`
	}{
		Name:       p.Name,
		Address:    p.Address.String(),
		Kind:       p.Kind,
		ID:         p.ID,
		APIAddress: p.APIAddress.String(),
		Secret:     p.Secret,
	}
	return json.Marshal(aux)
}
