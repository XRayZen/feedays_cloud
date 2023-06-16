package Data


type SearchResult struct{
	// refuse, accept
	ApiResponse string `json:"apiResponse"`
	ResponseMessage string `json:"responseMessage"`
	// found, none, error
	ResultType string `json:"resultType"`
	Exception string `json:"exception"`
	//AddContentならサイトを返す
    //PowerSearchなら記事を返す
	SearchType string `json:"searchType"`
	Websites []WebSite `json:"websites"`
	Articles []Article `json:"articles"`
}


