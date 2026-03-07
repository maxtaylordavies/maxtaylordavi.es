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
			Title:  "deployer",
			Status: "finished",
			Tags:   []string{"side", "web"},
			Link:   "https://github.com/maxtaylordavies/deployer",
		},
		{
			ID:     2,
			Title:  "switchboard",
			Status: "finished",
			Tags:   []string{"side", "web"},
			Link:   "https://github.com/maxtaylordavies/switchboard",
		},
		{
			ID:     3,
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

func (pr *ProjectRepository) WithTag(tag string) ([]Project, error) {
	var projects []Project

	allProjects, err := pr.All()
	if err != nil {
		return projects, err
	}

	for _, p := range allProjects {
		if contains(p.Tags, tag) {
			projects = append(projects, p)
		}
	}

	return projects, nil
}
