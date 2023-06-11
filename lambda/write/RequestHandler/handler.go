package RequestHandler

import (
	"errors"
	"write/DBRepo"
)

// ParseRequestType はリクエストタイプに応じて処理を分岐する
func ParseRequestType(dbRepo DBRepo.DBRepo, requestType string, userId string,
	argument_value1 string, argument_value2 string) (string, error) {
	functions := APIFunctions{
		repo: dbRepo,
	}
	switch requestType {
	case "GenUserID":
		return GenRandomUserID(argument_value1)
	case "CodeSync":
		return functions.CodeSync(argument_value1, argument_value2)
	case "GetUserInfo":
		// この機能はテスト用なので、実際には使わない
		return functions.GetUserInfo(userId)
	case "RegisterUser":
		return functions.RegisterUser(userId, argument_value1, argument_value2)
	case "ReportActivity":
		return functions.ReportActivity(argument_value1)
	case "SyncConfig":
		return functions.SyncConfig(argument_value1)
	case "editRecentSearches":
		return functions.editRecentSearches(userId, argument_value1, argument_value2)
	case "favoriteSite":
		return functions.favoriteSite(userId, argument_value1, argument_value2)
	case "favoriteArticle":
		return functions.favoriteArticle(userId, argument_value1, argument_value2)
	default:
		return "", errors.New("invalid request type")
	}
}
