package engine

var GraphQLJava = &engine{
	Name: "",
	Imprints: []imprint{
		{
			Query:   "",
			Matcher: inResponseText([]string{"Invalid Syntax : offending token '<EOF>'"}),
		},
		{
			Query:   "query @aaa@aaa { __typename }",
			Matcher: inResponseText([]string{"Validation error of type DuplicateDirectiveName: Directives must be uniquely named within a location."}),
		},
		{
			Query:   "queryy { __typename }",
			Matcher: inResponseText([]string{"Invalid Syntax : offending token 'queryy'"}),
		},
	},
}
