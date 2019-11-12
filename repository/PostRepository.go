package repository

import (
	"database/sql"
	"github.com/PuerkitoBio/goquery"
	"os"
	"path/filepath"
	"strconv"
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

	err := filepath.Walk("./posts", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || info.Name() == ".DS_Store" {
			return nil
		}

		r, err := os.Open(path)
		doc, err := goquery.NewDocumentFromReader(r)

		if err != nil {
			return err
		}

		str := info.Name()
		idStr := str[0]
		id, _ := strconv.Atoi(string(idStr))

		date, err := time.Parse("2006-01-02", doc.Find("p").Text()[:10])
		if err != nil {
			return err
		}

		posts = append(posts, Post{
			id,
			doc.Find("h1").Text(),
			doc.Find("p").Text()[10:],
			date,
		})

		return nil
	})

	return reversePostSlice(posts), err
}

func (pr * PostRepository) Recent() ([]Post, error) {
	var recentPosts []Post

	allPosts, err := pr.All()
	if err != nil {
		return recentPosts, err
	}

	if len(allPosts) > 3 {
		recentPosts = allPosts[:3]
	} else {
		recentPosts = allPosts
	}

	return recentPosts, nil
}

func (pr * PostRepository) One(id string) (Post, error) {
	var post Post

	err := pr.DatabaseConnection.QueryRow("SELECT * FROM posts where id = $1", id).Scan(&post.Id, &post.Title, &post.Body, &post.CreatedAt)
	if err != nil {
		return post, err
	}

	return post, nil
}

func reversePostSlice(slc []Post) []Post {
	for i, j := 0, len(slc)-1; i < j; i, j = i+1, j-1 {
		slc[i], slc[j] = slc[j], slc[i]
	}
	return slc
}