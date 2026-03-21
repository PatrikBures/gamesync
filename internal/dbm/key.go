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
	Comment string
}

func KeyAdd(db *sql.DB, pubKeyS string, user User) error {
	pk, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(pubKeyS))
	if err != nil {
		return fmt.Errorf("parsing pub key: %w", err)
	}
	fingerprint := ssh.FingerprintSHA256(pk)

	const SQL = `INSERT INTO ssh_key (user_id, fingerprint, pk, comment) VALUES(?, ?, ?, ?)`
	if _, err := db.Exec(SQL, user.ID, fingerprint, pk.Marshal(), comment); err != nil {
		return fmt.Errorf("inserting key: %w", err)
	}
	return nil
}

func KeyGetKeysForUserID(db *sql.DB, userID int) ([]Key, error) {
	const SQL = `SELECT key_id, user_id, fingerprint, pk, comment FROM ssh_key WHERE user_id = ?`
	rows, err := db.Query(SQL, userID)
	if err != nil {
		return nil, fmt.Errorf("selecting keys: %w", err)
	}
	var keys []Key
	for rows.Next() {
		k := Key{}
		if err := rows.Scan(&k.ID, &k.UserID, &k.Fingerprint, &k.PK, &k.Comment); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		keys = append(keys, k)
	}

	return keys, nil
}

func KeyGetKeysByFingerprint(db *sql.DB, fp string) ([]string, int, error) {
	const userIdSQL = `SELECT user_id FROM ssh_key WHERE fingerprint = ? LIMIT 1`

	var userID int
	row := db.QueryRow(userIdSQL, fp)
	if err := row.Scan(&userID); err != nil {
		return nil, -1, fmt.Errorf("selecting user id using fingerprint: %w", err)
	}

	const SQL = `SELECT pk FROM ssh_key where user_id = ?`
	rows, err := db.Query(SQL, userID)
	if err != nil {
		return nil, -1, fmt.Errorf("selecting keys: %w", err)
	}
	var keys []string
	for rows.Next() {
		var k []byte
		if err := rows.Scan(&k); err != nil {
			return nil, -1, fmt.Errorf("scanning row: %w", err)
		}
		pk, err := ssh.ParsePublicKey(k)
		if err != nil {
			return nil, -1, fmt.Errorf("parsing public key: %w", err)
		}
		key := ssh.MarshalAuthorizedKey(pk)
		keys = append(keys, string(key))
	}
	return keys, userID, nil
}
