package engine

var AWSAppSync = &engine{
	Name: "AWSAppSync",
	Imprints: []imprint{
		{
			Query:   "query @skip { __typename }",
			Matcher: inResponseText([]string{"MisplacedDirective"}),
		},
	},
}
