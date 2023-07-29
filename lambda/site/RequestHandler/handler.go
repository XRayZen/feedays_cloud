package RequestHandler

import (
	"site/Repo"
	"site/APIFunction"
)



func ParseRequestType(access_ip string, db_repo Repo.DBRepository, request_type string, user_id string, request_argument_json1 string, request_argument_json2 string) (string, error) {
	functions :=APIFunction.APIFunctions{
		DBRepo: db_repo,
	}
	switch request_type {
	case "Search":
		return functions.Search(access_ip, user_id, request_argument_json1)
	case "SubscribeSite":
		return functions.SubscribeSite(access_ip, user_id, request_argument_json1, request_argument_json2)
	case "fetchArticle":
		return functions.FetchArticle(access_ip, user_id, request_argument_json1)
	default:
		return "", nil
	}
}
