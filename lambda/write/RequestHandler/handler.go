package RequestHandler

import (
	"errors"
	"write/DBRepo"
)

// ParseRequestType はリクエストタイプに応じて処理を分岐する
func ParseRequestType(access_ip string, dbRepo DBRepo.DBRepo, requestType string, userId string,
	argumentJson_1 string, argumentJson_2 string) (string, error) {
	functions := APIFunctions{
		repo: dbRepo,
		ip:   access_ip,
	}
	switch requestType {
	case "GenUserID":
		return GenRandomUserID(dbRepo, argumentJson_1, access_ip)
	case "ConfigSync":
		return functions.ConfigSync(userId, argumentJson_1)
	case "GetUserInfo":
		// この機能はDB読み書きテスト用
		return functions.GetUserInfo(userId)
	case "RegisterUser":
		return functions.RegisterUser(userId, argumentJson_1, argumentJson_2)
	case "ReportReadActivity":
		return functions.ReportReadActivity(userId, argumentJson_1, argumentJson_2)
	case "UpdateConfig":
		return functions.UpdateConfig(userId, argumentJson_1, argumentJson_2)
	case "ModifySearchHistory":
		return functions.ModifySearchHistory(userId, argumentJson_1, argumentJson_2)
	case "ModifyFavoriteSite":
		return functions.ModifyFavoriteSite(userId, argumentJson_1, argumentJson_2)
	case "ModifyFavoriteArticle":
		return functions.ModifyFavoriteArticle(userId, argumentJson_1, argumentJson_2)
	case "GetAPIRequestLimit":
		return functions.GetAPIRequestLimit(userId, argumentJson_1)
	default:
		return "", errors.New("invalid request type")
	}
}
