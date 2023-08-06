package api

type Admin struct {
	AdminID  string `db:"admin_id"`
	Password string `db:"password"`
}
