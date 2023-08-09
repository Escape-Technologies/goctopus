package engine

var Agoo = &engine{
	Name: "agoo",
	Imprints: []imprint{
		{
			Query:   "query { zzz }",
			Matcher: inSection("code", []string{"eval error"}),
		},
	},
}
