package dbm

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type key struct {
	ID int
}

func KeyAdd(db *sql.DB, pubKeyS string, user User) error {
	pk, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(pubKeyS))
	if err != nil {
		return fmt.Errorf("parsing pub key: %w", err)
	}
	fingerprint := ssh.FingerprintSHA256(pk)

	const SQL = `INSERT INTO ssh_key (user_id, fingerprint, pk, type, comment) VALUES(?, ?, ?, ?, ?)`
	if _, err := db.Exec(SQL, user.ID, fingerprint, pk.Marshal(), pk.Type(), comment); err != nil {
		return fmt.Errorf("inserting key: %w", err)
	}
	return nil
}
