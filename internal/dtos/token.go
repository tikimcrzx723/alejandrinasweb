package dtos

type InsertToken struct {
	Plaintext string `json:"token"`
	Hash      []byte `json:"-"`
	UserID    string `json:"-"`
	Expiry    int64  `json:"expiry"`
	Scope     string `json:"-"`
}
