package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
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

	rows, err := pr.DatabaseConnection.Query(`SELECT * FROM projects`)
	if err != nil {
		return projects, err
	}

	for rows.Next() {
		p := Project{}
		if err := rows.Scan(&p.Id, &p.Title, &p.Body, &p.Status); err != nil {
			return projects, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (pr * ProjectRepository) Recent() ([]Project, error) {
	var projects []Project

	rows, err := pr.DatabaseConnection.Query(`SELECT * FROM posts ORDER BY id DESC limit 3`)
	if err != nil {
		return projects, err
	}

	for rows.Next() {
		p := Project{}
		if err := rows.Scan(&p.Id, &p.Title, &p.Body, &p.Status); err != nil {
			return projects, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (pr * ProjectRepository) One(id string) (Project, error) {
	var project Project

	err := pr.DatabaseConnection.QueryRow("SELECT * FROM projects where id = $1", id).Scan(&project.Id, &project.Title, &project.Body, &project.Status)
	if err != nil {
		return project, err
	}

	return project, nil
}


