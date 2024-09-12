package domain

type Player struct {
	ID         uint    `json:"id"`
	PlayerName string  `json:"player_name"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
}

type IPlayerService interface {
	CreatePlayers(players *[]Player) error
	GetPlayersGeoJSON() (string, error)
}

type IPlayerRepository interface {
	CreateBatch(*[]Player) error
	GetPlayers() (*[]Player, error)
}
