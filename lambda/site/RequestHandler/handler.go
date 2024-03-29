package RequestHandler

import (
	"fmt"
	"site/ApiFunction"
	"site/Repo"
)

// リクエストタイプはSearch, SubscribeSite, FetchArticleの3種類
func ParseRequestType(access_ip string, db_repo Repo.DBRepository, request_type string, user_id string, request_argument_json1 string, request_argument_json2 string) (string, error) {
	functions := ApiFunction.APIFunctions{
		DBRepo: db_repo,
	}
	switch request_type {
	case "Search":
		return functions.Search(access_ip, user_id, request_argument_json1)
	case "SubscribeSite":
		return functions.SubscribeSite(access_ip, user_id, request_argument_json1, request_argument_json2)
	case "FetchArticle":
		return functions.FetchArticle(access_ip, user_id, request_argument_json1)
	case "ModifyExploreCategory":
		return functions.ModifyExploreCategory(request_argument_json1, request_argument_json2)
	case "ChangeSiteCategory":
		return functions.ChangeSiteCategory(access_ip, user_id, request_argument_json1, request_argument_json2)
	case "DeleteSite":
		return functions.DeleteSite(access_ip, request_argument_json1, request_argument_json2)
	default:
		// リクエストタイプが不正な場合はエラーを返す
		return "", fmt.Errorf("invalid request type")
	}
}
