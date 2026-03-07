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
	CreatedAt time.Time
	Tags      []string
}

func (pr *PostRepository) All() ([]Post, error) {
	var posts []Post

	err := filepath.Walk("./posts", func(path string, info os.FileInfo, err error) error {
		fn := info.Name()

		if fn == "posts" || fn == "media" || strings.Contains(fn, ".") {
			return nil
		}

		// parse id
		idStr := fn[0]
		id, _ := strconv.Atoi(string(idStr))

		// get contents of main.html as string
		b, err := ioutil.ReadFile("./posts/" + fn + "/main.html")
		if err != nil {
			return err
		}
		s := string(b)

		// parse title
		i := strings.Index(s, "<title>")
		j := strings.Index(s, "</title>")
		title := strings.ToLower(s[i+7 : j])

		// parse date
		date, err := time.Parse("Mon Jan 02 2006", parseHtmlMetaContent(s, "date")[:15])
		if err != nil {
			return err
		}

		// parse tags
		tags := strings.Split(parseHtmlMetaContent(s, "tags"), ",")

		posts = append(posts, Post{
			id,
			title,
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

func (pr *PostRepository) WithTag(tag string) ([]Post, error) {
	var posts []Post

	allPosts, err := pr.All()
	if err != nil {
		return posts, err
	}

	for _, p := range allPosts {
		if contains(p.Tags, tag) {
			posts = append(posts, p)
		}
	}

	return posts, nil
}

func reversePostSlice(slc []Post) []Post {
	for i, j := 0, len(slc)-1; i < j; i, j = i+1, j-1 {
		slc[i], slc[j] = slc[j], slc[i]
	}
	return slc
}

func parseHtmlMetaContent(html string, key string) string {
	target := fmt.Sprintf(`<meta name="%s"`, key)
	i := strings.Index(html, target) + 23 + len(key)
	content := ""
	for {
		if string(html[i]) == `"` {
			break
		}
		content += string(html[i])
		i++
	}
	return content
}

func contains(slc []string, s string) bool {
	for _, a := range slc {
		if a == s {
			return true
		}
	}
	return false
}
