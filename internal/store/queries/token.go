package queries

const (
	TokenInsert = `
		INSERT INTO tokens(token, user_id, expiry)
		VALUES ($1, $2, $3)
		RETURNING id`
)
