package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	Id        int
	UserId    int
	TokenHash string
	ExpiresAt time.Time
	// Token is only set when a PasswordReset is being created
	Token string
}

// NOTE: we're setting the type of ExpiresAt to time.Time, because it's never
// going to be null. if it could be null, then we had to use sql.NullTime
// instead

//

type PasswordResetService struct {
	DB *sql.DB
	// how many bytes to use when generating a token. if it's not set, or the
	// provided value is less than MinBytesPerToken contant, then the default
	// value will be MinBytesPerToken.
	BytesPerToken int
	// the amount of time that a PasswordReset is valid
	Duration time.Duration
}

// TODO: implement this
func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: this function is not implemented")
}

// TODO: implement this
func (service *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: this function is not implemented")
}
