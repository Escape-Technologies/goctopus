package engine

var GraphQLGopherGo = &engine{
	Name: "",
	Imprints: []imprint{
		{
			Query:   "query {}",
			Matcher: hasJsonKey("data"),
		},
	},
}
