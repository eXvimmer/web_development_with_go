package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/exvimmer/lenslocked/rand"
)

const (
	// The minimum number of bytes to be used for each session token.
	MinBytePerToken = 32
)

type Session struct {
	Id        int
	UserId    int
	TokenHash string
	// Token is only set when creating a new session. When looking up a session,
	// this will be left empty, as we only store the hash of a session token in
	// our database and we cannot reverse it into a raw token.
	Token string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each session token. If this value is not set or is less than the
	// MinBytePerToken const, it will be ignored and MinBytePerToken will be
	// used.
	BytesPerToken int
}

func (ss *SessionService) Create(userId int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytePerToken {
		bytesPerToken = MinBytePerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("models/Create: %w", err)
	}
	session := &Session{
		Token:     token,
		UserId:    userId,
		TokenHash: ss.hash(token),
	}
	// TODO: store session in DB
	return session, nil
}

// TODO: needs implementation
func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}

func (ss *SessionService) hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(sum[:])
}
