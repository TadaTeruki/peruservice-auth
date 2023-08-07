package api

type Admin struct {
	AdminID  string `db:"id"`
	Password string `db:"password"`
}
