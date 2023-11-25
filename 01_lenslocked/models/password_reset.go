package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/exvimmer/lenslocked/rand"
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

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)
	var userId int
	row := service.DB.QueryRow(`
		SELECT id FROM users
		WHERE email = $1;
	`, email)
	err := row.Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytePerToken {
		bytesPerToken = MinBytePerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	duration := service.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}
	pwReset := PasswordReset{
		UserId:    userId,
		Token:     token,
		TokenHash: service.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}
	row = service.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES
			($1, $2, $3)
		ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;
	`, pwReset.UserId, pwReset.TokenHash, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.Id)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &pwReset, nil
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := service.hash(token)
	row := service.DB.QueryRow(`
		SELECT password_resets.id, password_resets.expires_at,
		users.id, users.email, users.password_hash
		FROM password_resets
		JOIN users
		ON users.id = password_resets.user_id
		WHERE password_resets.token_hash = $1;
	`, tokenHash)
	var (
		user    User
		pwReset PasswordReset
	)
	err := row.Scan(
		&pwReset.Id,
		&pwReset.ExpiresAt,
		&user.Id,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}
	err = service.delete(pwReset.Id)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	return &user, nil
}

func (service *PasswordResetService) hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(sum[:])
}

func (service *PasswordResetService) delete(id int) error {
	_, err := service.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
