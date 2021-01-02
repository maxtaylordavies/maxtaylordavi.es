package repository

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ProjectRepository struct {
}

type Project struct {
	Id     int
	Title  string
	Status string
	Tags   []string
}

func (pr *ProjectRepository) All() ([]Project, error) {
	var projects []Project

	err := filepath.Walk("./projects", func(path string, info os.FileInfo, err error) error {
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

		lines := strings.Split(doc.Find("p").Text(), "\n")

		var tags []string
		if len(lines[0]) > 21 {
			tags = strings.Split(lines[0][22:], " ")
		}

		projects = append(projects, Project{
			id,
			strings.ToLower(doc.Find("h1").Text()),
			"Done",
			tags,
		})

		return nil
	})

	return reverseSlice(projects), err
}

func (pr *ProjectRepository) Recent() ([]Project, error) {
	var recentProjects []Project

	allProjects, err := pr.All()
	if err != nil {
		return recentProjects, err
	}

	if len(allProjects) > 3 {
		recentProjects = allProjects[:3]
	} else {
		recentProjects = allProjects
	}

	return recentProjects, nil
}

func reverseSlice(slc []Project) []Project {
	for i, j := 0, len(slc)-1; i < j; i, j = i+1, j-1 {
		slc[i], slc[j] = slc[j], slc[i]
	}
	return slc
}
