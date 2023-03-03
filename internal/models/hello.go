package models

// Hello represents message from one person to another
type Hello struct {
	ID      int64  `db:"id"`
	From    string `db:"from"`
	To      string `db:"to"`
	Message string `db:"message"`
}
