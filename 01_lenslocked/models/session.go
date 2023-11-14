package models

import "database/sql"

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

// TODO: needs implementation
func (ss *SessionService) Create(userId int) (*Session, error) {
	return nil, nil
}

// TODO: needs implementation
func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}
