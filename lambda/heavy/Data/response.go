package Data


type SearchResult struct{
	// refuse, accept
	ApiResponse string
	ResponseMessage string
	// found, none, error
	ResultType string
	Exception string
	//AddContentならサイトを返す
    //PowerSearchなら記事を返す
	SearchType string
	Websites []WebSite
	Articles []Article
}


