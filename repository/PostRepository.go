package repository

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type PostRepository struct {
}

type Post struct {
	Id        int
	Title     string
	Body      string
	CreatedAt time.Time
	Tags      []string
}

func (pr *PostRepository) All() ([]Post, error) {
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

		pText := strings.Split(doc.Find("p").Text(), "\n")[0]

		date, err := time.Parse("2006-01-02", pText[:10])
		if err != nil {
			return err
		}

		var tags []string
		if len(pText) > 10 {
			tags = strings.Split(pText[11:], " ")
		}

		posts = append(posts, Post{
			id,
			strings.ToLower(doc.Find("h1").Text()),
			doc.Find("p").Text()[10:],
			date,
			tags,
		})

		return nil
	})

	return reversePostSlice(posts), err
}

func (pr *PostRepository) Recent() ([]Post, error) {
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

func reversePostSlice(slc []Post) []Post {
	for i, j := 0, len(slc)-1; i < j; i, j = i+1, j-1 {
		slc[i], slc[j] = slc[j], slc[i]
	}
	return slc
}
