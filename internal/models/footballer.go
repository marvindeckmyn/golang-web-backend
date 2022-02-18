package models

type Footballer struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FootballClub string `json:"football_club"`
}
