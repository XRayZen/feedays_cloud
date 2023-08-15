package RequestHandler

import (
	"errors"
	"user/DBRepo"
)

// ParseRequestType はリクエストタイプに応じて処理を分岐する
func ParseRequestType(access_ip string, db_repo DBRepo.DBRepo, request_type string, user_id string,
	argument_json_1 string, argument_json_2 string) (string, error) {
	functions := APIFunctions{
		db_repo: db_repo,
		ip:      access_ip,
	}
	// リクエストタイプに応じて処理を分岐
	switch request_type {
	case "ServiceInitialize":
		return functions.ServiceInitialize()
	case "ServiceFinalize":
		return functions.ServiceFinalize()
	case "GenUserID":
		return GenRandomUserID()
	case "RegisterUser":
		return functions.RegisterUser(user_id, argument_json_1)
	case "ConfigSync":
		return functions.ConfigSync(user_id)
	case "ReportReadActivity":
		return functions.ReportReadActivity(user_id, argument_json_1)
	case "UpdateConfig":
		return functions.UpdateConfig(user_id, argument_json_1)
	case "ModifySearchHistory":
		return functions.ModifySearchHistory(user_id, argument_json_1, argument_json_2)
	case "ModifyFavoriteSite":
		return functions.ModifyFavoriteSite(user_id, argument_json_1, argument_json_2)
	case "ModifyFavoriteArticle":
		return functions.ModifyFavoriteArticle(user_id, argument_json_1, argument_json_2)
	case "GetAPIRequestLimit":
		return functions.GetAPIRequestLimit(user_id)
	case "ModifyAPIRequestLimit":
		return functions.ModifyAPIRequestLimit(argument_json_1, argument_json_2)
	case "DeleteUserData":
		return functions.DeleteUserData(user_id, argument_json_1)
	default:
		return "", errors.New("invalid request type")
	}
}
