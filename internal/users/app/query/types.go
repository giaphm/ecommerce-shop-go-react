package query

type UserModel struct {
	uuid           string
	DisplayName    string
	Email          string
	HashedPassword string
	Balance        int
	Role           string
	LastIP         string
}
