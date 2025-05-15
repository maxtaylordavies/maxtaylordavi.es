package repository

type PaperRepository struct {
}

type Paper struct {
	ID      int
	Title   string
	URL     string
	Authors string
	Venue   string
	Tags    []string
}

type PaperBatch struct {
	Key    string
	Papers []Paper
}

func reverse(in []Paper) []Paper {
	out := make([]Paper, len(in))
	for i, p := range in {
		out[len(in)-1-i] = p
	}
	return out
}

func (pr *PaperRepository) All() ([]PaperBatch, error) {
	data := []PaperBatch{
		{
			Key: "Preprints",
			Papers: reverse([]Paper{
				{
					ID:      7,
					Title:   "Factive mindreading reflects the optimal use of limited cognitive resources",
					URL:     "https://osf.io/preprints/psyarxiv/wzbcs_v1",
					Authors: "Tadeg Quillien + Max Taylor-Davies",
					Venue:   "",
					Tags:    []string{},
				},
			}),
		},
		{
			Key: "2025",
			Papers: reverse([]Paper{
				{
					ID:      4,
					Title:   "Emergent kin selection of altruistic feeding via non-episodic neuroevolution",
					URL:     "https://link.springer.com/chapter/10.1007/978-3-031-90062-4_31",
					Authors: "Max Taylor-Davies, Gautier Hamon, Timothé Boulet + Clément Moulin-Frier",
					Venue:   "28th International Conference on the Applications of Evolutionary Computation (✨best paper✨)",
					Tags:    []string{"conference", "talk", "award"},
				},
				{
					ID:      5,
					Title:   "An information bottleneck view of social stereotype use",
					URL:     "https://osf.io/preprints/psyarxiv/k53tz_v1",
					Authors: "Max Taylor-Davies + Tadeg quillien",
					Venue:   "Proceedings of the Annual Meeting of the Cognitive Science Society 2025",
					Tags:    []string{"conference", "talk"},
				},
				{
					ID:      6,
					Title:   "A rational framework for group-based selective social learning",
					URL:     "",
					Authors: "Max Taylor Davies, Neil Bramley + Christopher G. Lucas",
					Venue:   "Open Mind (forthcoming)",
					Tags:    []string{"journal"},
				},
			}),
		},
		{
			Key: "2024",
			Papers: reverse([]Paper{
				{
					ID:      3,
					Title:   "Rational compression in choice prediction",
					URL:     "/files/iccm-paper.pdf",
					Authors: "Max Taylor Davies + Christopher G. Lucas",
					Venue:   "22nd International Conference on Cognitive Modeling (ICCM), 2024",
					Tags:    []string{"conference", "talk"},
				},
			}),
		},
		{
			Key: "2023",
			Papers: reverse([]Paper{
				{
					ID:      0,
					Title:   "Selective imitation on the basis of reward function similarity",
					URL:     "https://escholarship.org/uc/item/8x3072nr",
					Authors: "Max Taylor Davies, Stephanie Droop + Christopher G. Lucas",
					Venue:   "Proceedings of the Annual Meeting of the Cognitive Science Society 2023",
					Tags:    []string{"conference", "poster"},
				},
				{
					ID:      1,
					Title:   "Is feedback all you need? leveraging natural language feedback in goal-conditioned reinforcement learning",
					URL:     "https://arxiv.org/abs/2312.04736",
					Authors: "Sabrina McCallum, Max Taylor Davies, Stefano Albrecht + Alessandro Suglia",
					Venue:   "Goal-Conditioned Reinforcement Learning Workshop, NeurIPS 2023 (spotlight)",
					Tags:    []string{"workshop", "talk"},
				},
				{
					ID:      2,
					Title:   "Balancing utility and cognitive cost in social representation",
					URL:     "https://arxiv.org/abs/2310.04852",
					Authors: "Max Taylor Davies + Christopher G. Lucas",
					Venue:   "Information-Theoretic Principles in Cognitive Systems Workshop, NeurIPS 2023",
					Tags:    []string{"workshop", "poster"},
				},
			}),
		},
	}

	return data, nil
}

func (pr *PaperRepository) WithTag(tag string) ([]PaperBatch, error) {
	var data []PaperBatch

	all, err := pr.All()
	if err != nil {
		return data, err
	}

	for _, batch := range all {
		var papers []Paper
		for _, p := range batch.Papers {
			if contains(p.Tags, tag) {
				papers = append(papers, p)
			}
		}
		if len(papers) > 0 {
			data = append(data, PaperBatch{
				Key:    batch.Key,
				Papers: papers,
			})
		}
	}

	return data, nil
}
