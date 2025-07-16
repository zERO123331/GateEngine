package data

type Player struct {
	Name       string
	UUID       string
	Rank       string
	systemRank string
	Banned     bool
}

func (p *Player) SetSystemRank(rank string) {
	if ValidateSystemRank(rank) {
		p.systemRank = rank
	}

}

func (p *Player) GetSystemRank() string {
	return p.systemRank
}

func (p *Player) Ban() {
	p.Banned = true
}

func (p *Player) Unban() {
	p.Banned = false
}

func ValidateSystemRank(rank string) bool {
	return rank == "admin" || rank == "mod" || rank == "player"
}
