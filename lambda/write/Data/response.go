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

type ConfigSyncResponse struct {
	ResponseType string     `json:"responseType"`
	UserConfig   UserConfig `json:"userConfig"`
	Error        string     `json:"error"`
}
