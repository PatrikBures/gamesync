package service

import "crypto/rand"

func generateToken() ([]byte, error) {
	b := make([]byte, 33)
	_, err := rand.Read(b)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
