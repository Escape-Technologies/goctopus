package engine

var Graphene = &engine{
	Name: "Graphene",
	Imprints: []imprint{
		{
			Query:   "query { aaa }",
			Matcher: inResponseText([]string{"Syntax Error GraphQL (1:1)"}),
		},
	},
}
