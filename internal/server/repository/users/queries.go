package usersrepo

const (
	CreateUserQuery = `INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id`
	GetUserIDQuery  = `SELECT id FROM users WHERE username = $1 AND password = $2`
)
