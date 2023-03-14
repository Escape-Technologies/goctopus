package engine

var GraphQLPHP = &engine{
	Name: "GraphQLPHP",
	Imprints: []imprint{
		{
			Query:   "query ! {__typename}",
			Matcher: inResponseText([]string{"Syntax Error: Cannot parse the unexpected character \"?\"."}),
		},
		{
			Query:   "query @deprecated {__typename}",
			Matcher: inResponseText([]string{"Directive \"deprecated\" may not be used on \"QUERY\"."}),
		},
	},
}
