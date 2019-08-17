package repository

import (
	"database/sql"
	"time"
)

type PostRepository struct {
	DatabaseConnection *sql.DB
}

type Post struct {
	Id int
	Title string
	Body string
	CreatedAt time.Time
}

func (pr * PostRepository) All() ([]Post, error) {
	var posts []Post

	rows, err := pr.DatabaseConnection.Query(`SELECT * FROM posts`)
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		p := Post{}
		if err := rows.Scan(&p.Id, &p.Title, &p.Body, &p.CreatedAt); err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (pr * PostRepository) Recent() ([]Post, error) {
	var posts []Post

	rows, err := pr.DatabaseConnection.Query(`SELECT * FROM posts ORDER BY id DESC limit 3`)
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		p := Post{}
		if err := rows.Scan(&p.Id, &p.Title, &p.Body, &p.CreatedAt); err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (pr * PostRepository) One(id string) (Post, error) {
	var post Post

	err := pr.DatabaseConnection.QueryRow("SELECT * FROM posts where id = $1", id).Scan(&post.Id, &post.Title, &post.Body, &post.CreatedAt)
	if err != nil {
		return post, err
	}

	return post, nil
}
