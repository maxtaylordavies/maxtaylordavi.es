package repository

import (
	"database/sql"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type ProjectRepository struct {
	DatabaseConnection *sql.DB
}

type Project struct {
	Id int
	Title string
	Body string
	Status string
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
			log.Println("a")
			return err
		}

		str := info.Name()
		idStr := str[0]
		id, _ := strconv.Atoi(string(idStr))

		projects = append(projects, Project{
			id,
			doc.Find("h1").Text(),
			doc.Find("p").Text()[10:],
			"Done",
		})

		return nil
	})

	return reverseSlice(projects), err
}

func (pr * ProjectRepository) Recent() ([]Project, error) {
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

func (pr * ProjectRepository) One(id string) (Project, error) {
	var project Project

	err := pr.DatabaseConnection.QueryRow("SELECT * FROM projects where id = $1", id).Scan(&project.Id, &project.Title, &project.Body, &project.Status)
	if err != nil {
		return project, err
	}

	return project, nil
}

func reverseSlice(slc []Project) []Project {
	for i, j := 0, len(slc)-1; i < j; i, j = i+1, j-1 {
		slc[i], slc[j] = slc[j], slc[i]
	}
	return slc
}


