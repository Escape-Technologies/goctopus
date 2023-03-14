package engine

var T = &engine{
	Name: "",
	Imprints: []imprint{
		{
			Query:   "",
			Matcher: inResponseText([]string{""}),
		},
	},
}
