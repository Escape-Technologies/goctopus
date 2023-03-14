package engine

var Apollo = &engine{
	Name: "apollo",
	Imprints: []imprint{
		{
			Query: "query @deprecated { __typename }",
			Matcher: inResponseText([]string{
				"Directive \\\"@deprecated\\\" may not be used on QUERY.",
				"Directive \\\"deprecated\\\" may not be used on QUERY.",
			}),
		},
		{
			Query: "query @skip { __typename }",
			Matcher: inResponseText([]string{
				"Directive \\\"@skip\\\" argument \\\"if\\\" of type \\\"Boolean!\\\" is required, but it was not provided",
				"Directive \\\"skip\\\" argument \\\"if\\\" of type \\\"Boolean!\\\" is required, but it was not provided",
			}),
		},
	},
}
