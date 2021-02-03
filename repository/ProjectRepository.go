package repository

type ProjectRepository struct {
}

type Project struct {
	ID     int
	Title  string
	Status string
	Tags   []string
	Link   string
}

func (pr *ProjectRepository) All() ([]Project, error) {
	projects := []Project{
		{
			ID:     1,
			Title:  "products.gallery",
			Status: "finished",
			Tags:   []string{"side", "web"},
			Link:   "http://products.gallery",
		},
	}

	return projects, nil
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
