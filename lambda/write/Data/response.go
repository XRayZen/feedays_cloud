package Data

type APIResponse struct {
	ResponseType string `json:"responseType"`
	Message      string `json:"message"`
	Error        string `json:"error"`
}

type FetchCloudFeedResponse struct {
	ResponseType string    `json:"responseType"`
	Feeds        []Article `json:"feeds"`
	Error        string    `json:"error"`
}

type CodeSyncResponse struct {
	ResponseType string     `json:"responseType"`
	UserConfig   UserConfig `json:"userConfig"`
	WebSites     []WebSite  `json:"webSites"`
	Error        string     `json:"error"`
}
