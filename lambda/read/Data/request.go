package Data

type ApiSearchRequest struct{
	// addContent, exploreWeb, powerSearch
	SearchType string `json:"searchType"`
	Word string `json:"word"`
	UserID string `json:"userID"`
	IdentInfo string `json:"identInfo"`
	AccountType string `json:"accountType"`
	RequestTime string `json:"requestTime"`
}

