package engine

var GraphQLGo = &engine{
	Name: "GraphQLGo",
	Imprints: []imprint{
		{
			Query:   "",
			Matcher: inResponseText([]string{"Must provide an operation."}),
		},
		{
			Query:   "query  { __typename {}",
			Matcher: inResponseText([]string{"Unexpected empty IN"}),
		},
		{
			Query:   "query  { __typename }",
			Matcher: inResponseText([]string{"RootQuery"}),
		},
	},
}
