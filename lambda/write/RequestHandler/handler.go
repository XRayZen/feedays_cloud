package RequestHandler

import (
	"errors"
	"write/DBRepo"
)

// ParseRequestType はリクエストタイプに応じて処理を分岐する
func ParseRequestType(repo DBRepo.DBRepo, requestType string, userId string) (interface{}, error) {
	switch requestType {
	case "GetUserInfo":
		return RequestType.GetUserInfo(repo, userId)
	case "GetUserList":
		return RequestType.GetUserList(repo)
	default:
		return nil, errors.New("invalid request type")
	}
}