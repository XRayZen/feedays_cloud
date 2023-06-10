package RequestHandler

import (
	"errors"
	"write/DBRepo"
)

// ParseRequestType はリクエストタイプに応じて処理を分岐する
func ParseRequestType(repo DBRepo.DBRepo, requestType string, userId string,
	argument_value1 string, argument_value2 string) (string, error) {
	switch requestType {
	case "GenUserID":
		return GenRandomUserID(argument_value1)
	case "CodeSync":
		return CodeSync(repo, userId, argument_value1)
	case "GetUserInfo":
		// この機能はテスト用なので、実際には使わない
		return GetUserInfo(repo, userId)
	case "RegisterUser":
		return RegisterUser(repo, userId, argument_value1, argument_value2)
	case "ReportActivity":
		return ReportActivity(repo, userId, argument_value1)
	case "SyncConfig":
		return SyncConfig(repo, userId, argument_value1)
	default:
		return "", errors.New("invalid request type")
	}
}
