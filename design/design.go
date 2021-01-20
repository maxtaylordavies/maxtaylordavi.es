package design

type Theme = struct {
	Background string
	Color      string
}

func GetTheme(name string) Theme {
	themes := map[string]Theme{
		"purple": {
			Background: "rgba(145, 71, 255, 0.2)",
			Color:      "rgb(145, 71, 255)",
		},
	}

	if theme, ok := themes[name]; ok {
		return theme
	}

	return themes["purple"]
}
