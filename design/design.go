package design

type Theme = struct {
	Name       string
	Background string
	Color      string
}

func GetAllThemes() []Theme {
	return []Theme{
		{
			Name:       "teal",
			Background: "rgba(152, 255, 212, 0.15)",
			Color:      "#50E0A4",
		},
		{
			Name:       "purple",
			Background: "rgba(145, 71, 255, 0.15)",
			Color:      "rgb(145, 71, 255)",
		},
		{
			Name:       "yellow",
			Background: "rgba(255, 175, 0, 0.15)",
			Color:      "rgb(255,175,0)",
		},
		{
			Name:       "blue",
			Background: "rgba(72, 88, 234, 0.15)",
			Color:      "rgb(72, 88, 234)",
		},
		{
			Name:       "coral",
			Background: "rgba(255, 156, 123, 0.15)",
			Color:      "rgb(255, 156, 123)",
		},
		{
			Name:       "pink",
			Background: "rgba(242, 114, 215, 0.15)",
			Color:      "rgb(242, 114, 215)",
		},
	}
}

func GetTheme(name string) Theme {
	themes := GetAllThemes()

	for _, theme := range themes {
		if theme.Name == name {
			return theme
		}
	}

	return themes[4]
}
