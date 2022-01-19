package model

type Room struct {
	ID           int
	Clients      map[string]*Client
	Password     string
	CurrentVideo *Video
}
