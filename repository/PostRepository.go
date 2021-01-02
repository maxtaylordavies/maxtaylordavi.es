package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

		// parse id
		fn := info.Name()
		idStr := fn[0]
		id, _ := strconv.Atoi(string(idStr))

		b, err := ioutil.ReadFile("./posts/" + fn)
		if err != nil {
			return err
		}
		s := string(b)

		// parse title
		i := strings.Index(s, "<h1>")
		j := strings.Index(s, "</h1>")
		title := strings.ToLower(s[i+23 : j])

		// parse body
		i = strings.Index(s, "</p><p")
		tempS := s[i+6:]
		i = strings.Index(tempS, ">")
		j = strings.Index(tempS, "</p")
		body := tempS[i+1 : j]

		// parse date
		i = strings.Index(s, "<em>")
		date, err := time.Parse("2006-01-02", s[i+4:i+14])
		if err != nil {
			return err
		}

		// parse tags
		var tags []string
		i = strings.Index(s, "name='tags'")
		if i != -1 {
			tempS = s[i:]
			j := strings.Index(tempS, "'/>")
			tags = strings.Split(tempS[21:j], " ")
		}

		fmt.Println(tags)

		posts = append(posts, Post{
			id,
			title,
			body,
			date,
			tags,
		})

		return nil
	})

	fmt.Println(posts)
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
