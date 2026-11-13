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
			Key: "Working papers",
			Papers: reverse([]Paper{
				{
					ID:      11,
					Title:   "Aligning Incentives in Cumulative Cultural Evolution",
					URL:     "/files/cce-working-paper.pdf",
					Authors: "Me, Anil Yaman",
					Venue:   "",
					Tags:    []string{},
				},
			}),
		},
		{
			Key: "Preprints",
			Papers: reverse([]Paper{
				{
					ID:      10,
					Title:   "Learning When to Cooperate Under Heterogeneous Goals",
					URL:     "https://arxiv.org/abs/2603.07253",
					Authors: "Me, Neil R. Bramley, Christopher G. Lucas",
					Venue:   "",
					Tags:    []string{},
				},
			}),
		},
		{
			Key: "2026",
			Papers: reverse([]Paper{
				{
					ID:      8,
					Title:   "Factive Mindreading Reflects the Optimal Use of Limited Cognitive Resources",
					URL:     "https://osf.io/preprints/psyarxiv/wzbcs_v1",
					Authors: "Tadeg Quillien, Me",
					Venue:   "Proceedings of the Royal Society B",
					Tags:    []string{"journal"},
				},
				{
					ID:      9,
					Title:   "Representational Capacity Shapes Both the \"Whether\" and \"How\" of Social Learning",
					URL:     "https://journals.plos.org/ploscompbiol/article?id=10.1371/journal.pcbi.1014504",
					Authors: "Me",
					Venue:   "PLOS Computational Biology",
					Tags:    []string{"journal"},
				},
			}),
		},
		{
			Key: "2025",
			Papers: reverse([]Paper{
				{
					ID:      4,
					Title:   "Emergent Kin Selection of Altruistic Feeding via Non-episodic Neuroevolution",
					URL:     "https://link.springer.com/chapter/10.1007/978-3-031-90062-4_31",
					Authors: "Me, Gautier Hamon, Timothé Boulet, Clément Moulin-Frier",
					Venue:   "28th International Conference on the Applications of Evolutionary Computation (✨best paper✨)",
					Tags:    []string{"conference", "talk", "award"},
				},
				{
					ID:      5,
					Title:   "An Information Bottleneck View of Social Stereotype Use",
					URL:     "https://osf.io/preprints/psyarxiv/k53tz_v1",
					Authors: "Me, Tadeg quillien",
					Venue:   "Proceedings of the Annual Meeting of the Cognitive Science Society 2025",
					Tags:    []string{"conference", "talk"},
				},
				{
					ID:      6,
					Title:   "A Rational Framework for Group-Based Selective Social Learning",
					URL:     "https://direct.mit.edu/opmi/article/doi/10.1162/opmi_a_00205/131045/A-Rational-Framework-for-Group-Based-Selective",
					Authors: "Me, Neil R. Bramley, Christopher G. Lucas",
					Venue:   "Open Mind",
					Tags:    []string{"journal"},
				},
				{
					ID:      7,
					Title:   "Partner Modelling Emerges in Recurrent Agents (But Only When it Matters)",
					URL:     "https://arxiv.org/abs/2505.17323",
					Authors: "Ruaridh Mon-Williams*, Me*, Elizabeth Mieczkowski, Natalia Velez, Neil R. Bramley, Yanwei Wang, Thomas L. Griffiths, Christopher G. Lucas<br/>(* co-first authors)",
					Venue:   "39th Annual Conference on Neural Information Processing Systems",
					Tags:    []string{"conference", "poster"},
				},
			}),
		},
		{
			Key: "2024",
			Papers: reverse([]Paper{
				{
					ID:      3,
					Title:   "Rational Compression in Choice Prediction",
					URL:     "/files/iccm-paper.pdf",
					Authors: "Me, Christopher G. Lucas",
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
					Title:   "Selective Imitation on the Basis of Reward Function Similarity",
					URL:     "https://escholarship.org/uc/item/8x3072nr",
					Authors: "Me, Stephanie Droop, Christopher G. Lucas",
					Venue:   "Proceedings of the Annual Meeting of the Cognitive Science Society 2023",
					Tags:    []string{"conference", "poster"},
				},
				{
					ID:      1,
					Title:   "Is Feedback All You Need? Leveraging Natural Language Feedback in Goal-Conditioned Reinforcement Learning",
					URL:     "https://arxiv.org/abs/2312.04736",
					Authors: "Sabrina McCallum, Me, Stefano Albrecht, Alessandro Suglia",
					Venue:   "Goal-Conditioned Reinforcement Learning Workshop, NeurIPS 2023 (spotlight)",
					Tags:    []string{"workshop", "talk"},
				},
				{
					ID:      2,
					Title:   "Balancing Utility and Cognitive Cost in Social Representation",
					URL:     "https://arxiv.org/abs/2310.04852",
					Authors: "Me, Christopher G. Lucas",
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
