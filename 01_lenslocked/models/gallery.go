package models

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type Image struct {
	GalleryId int
	Path      string
	Filename  string
}

type Gallery struct {
	Id     int
	Title  string
	UserId int
}

type GalleryService struct {
	DB *sql.DB
	// ImagesDir is used to tell the GalleryService where to store and locate
	// images. If not set, the GalleryService will default to using the
	// "images" directory.
	ImagesDir string
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

func (gs *GalleryService) ById(id int) (*Gallery, error) {
	gallery := Gallery{Id: id}
	row := gs.DB.QueryRow(`
		SELECT title, user_id FROM galleries
		WHERE id = $1;
	`, id)
	err := row.Scan(&gallery.Title, &gallery.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) ByUserId(userId int) ([]Gallery, error) {
	rows, err := gs.DB.Query(`
		SELECT id, title
		FROM galleries
		WHERE user_id = $1;
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{UserId: userId}
		err = rows.Scan(&gallery.Id, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user id: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user id: %w", err)
	}
	return galleries, nil
}

func (gs *GalleryService) Update(gallery *Gallery) error {
	if gallery == nil {
		return fmt.Errorf("gallery is nil")
	}
	_, err := gs.DB.Exec(`
		UPDATE galleries
		SET title = $1
		WHERE id = $2;
	`, gallery.Title, gallery.Id)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (gs *GalleryService) Delete(id int) error {
	_, err := gs.DB.Exec(`
		DELETE FROM galleries
		WHERE id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (gs *GalleryService) Images(galleryId int) ([]Image, error) {
	globPattern := filepath.Join(gs.galleryDir(galleryId), "*")
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving images: %w", err)
	}
	var images []Image
	for _, file := range matches {
		if hasExtension(file, gs.extensions()...) {
			images = append(images, Image{
				GalleryId: galleryId,
				Filename:  filepath.Base(file),
				Path:      file,
			})
		}
	}
	return images, nil
}

func (gs *GalleryService) galleryDir(id int) string {
	imagesDir := gs.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func (gs *GalleryService) extensions() []string {
	return []string{".jpg", ".jpeg", ".webp", ".gif", ".png"}
}

func hasExtension(file string, extensions ...string) bool {
	file = strings.ToLower(file)
	for _, ext := range extensions {
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}
