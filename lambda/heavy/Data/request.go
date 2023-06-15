package Data

type ApiSearchRequest struct{
	// addContent, exploreWeb, powerSearch
	SearchType string
	Word string
	UserID string
	IdentInfo string
	AccountType string
	RequestTime string
}

