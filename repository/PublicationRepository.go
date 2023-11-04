package repository

type PublicationRepository struct {
}

type Publication struct {
	ID      int
	Title   string
	URL     string
	Authors string
	Venue   string
	Tags    []string
}

type PublicationYear struct {
	Year int
	Pubs []Publication
}

func reverse(in []Publication) []Publication {
	out := make([]Publication, len(in))
	for i, p := range in {
		out[len(in)-1-i] = p
	}
	return out
}

func (pr *PublicationRepository) All() ([]PublicationYear, error) {
	data := []PublicationYear{
		{
			Year: 2023,
			Pubs: reverse([]Publication{
				{
					ID:      1,
					Title:   "Selective imitation on the basis of reward function similarity",
					URL:     "https://arxiv.org/abs/2305.07421",
					Authors: "Max Taylor Davies, Stephanie Droop + Christopher G. Lucas",
					Venue:   "Proceedings of the Annual Meeting of the Cognitive Science Society 2023",
					Tags:    []string{"conference"},
				},
				{
					ID:      2,
					Title:   "Is feedback all you need? leveraging natural language feedback in goal-conditioned reinforcement learning",
					URL:     "",
					Authors: "Sabrina McCallum, Max Taylor Davies, Stefano Albrecht + Alessandro Suglia",
					Venue:   "Goal-Conditioned Reinforcement Learning Workshop, NeurIPS 2023",
					Tags:    []string{"workshop", "spotlight"},
				},
				{
					ID:      3,
					Title:   "Balancing utility and cognitive cost in social representation",
					URL:     "https://arxiv.org/abs/2310.04852",
					Authors: "Max Taylor Davies + Christopher G. Lucas",
					Venue:   "Information-Theoretic Principles in Cognitive Systems Workshop, NeurIPS 2023",
					Tags:    []string{"workshop"},
				},
			}),
		},
	}

	return data, nil
}

func (pr *PublicationRepository) WithTag(tag string) ([]PublicationYear, error) {
	var data []PublicationYear

	all, err := pr.All()
	if err != nil {
		return data, err
	}

	for _, y := range all {
		var pubs []Publication
		for _, p := range y.Pubs {
			if contains(p.Tags, tag) {
				pubs = append(pubs, p)
			}
		}
		if len(pubs) > 0 {
			data = append(data, PublicationYear{
				Year: y.Year,
				Pubs: pubs,
			})
		}
	}

	return data, nil
}
