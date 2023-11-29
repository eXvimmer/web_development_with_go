package models

import (
	"database/sql"
	"fmt"
)

type Gallery struct {
	Id     int
	Title  string
	UserId int
}

type GalleryService struct {
	DB *sql.DB
}

func (gs *GalleryService) Create(title string, userId int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserId: userId,
	}
	// TODO: add validation
	row := gs.DB.QueryRow(`
		INSERT INTO galleries (title, user_id)
		VALUES ($1, $2)
		RETURNING id;
	`, title, userId)
	err := row.Scan(&gallery.Id)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}
