package models

import (
	"database/sql"
	"fmt"

	"github.com/exvimmer/lenslocked/rand"
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
}

func (ss *SessionService) Create(userId int) (*Session, error) {
	token, err := rand.SessionToken()
	if err != nil {
		return nil, fmt.Errorf("models/Create: %w", err)
	}
	session := &Session{
		Token:  token,
		UserId: userId,
		// TODO: hash the session token and set it here
	}
	// TODO: store session in DB
	return session, nil
}

// TODO: needs implementation
func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}
