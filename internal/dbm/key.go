package dbm

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type Key struct {
	ID int
	UserID int
	Fingerprint string
	PK []byte
	Type string
	Comment string
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

func KeyGetKeysForUserID(db *sql.DB, userID int) ([]Key, error) {
	const SQL = `SELECT key_id, user_id, fingerprint, pk, type, comment FROM ssh_key WHERE user_id = ?`
	rows, err := db.Query(SQL, userID)
	if err != nil {
		return nil, fmt.Errorf("selecting keys: %w", err)
	}
	var keys []Key
	for rows.Next() {
		k := Key{}
		if err := rows.Scan(&k.ID, &k.UserID, &k.Fingerprint, &k.PK, &k.Type, &k.Comment); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		keys = append(keys, k)
	}

	return keys, nil
}
