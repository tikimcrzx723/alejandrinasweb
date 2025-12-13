package queries

const (
	UserInsert = `
		INSERT INTO users(first_name, last_name, username, email, password_hash, activated, is_block)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	UserRoleInsert = `
		INSERT INTO users_roles (user_id, role_id)
	    VALUES ($1, $2)`
)
