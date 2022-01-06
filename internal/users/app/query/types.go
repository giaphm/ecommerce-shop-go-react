package query

type User struct {
	Uuid           string
	DisplayName    string
	Email          string
	HashedPassword string
	Balance        float32
	Role           string
	LastIP         string
}
