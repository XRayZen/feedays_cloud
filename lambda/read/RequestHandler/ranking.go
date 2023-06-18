package RequestHandler

import (
	// "errors"
	"encoding/json"
	"read/Repo"
)

type Ranking struct {
	DBrepo Repo.DBRepository
}

func (e Ranking) GetRanking(userID string) (RankingJson string, err error) {
	//DBからデータを取得する
	resUserInfo, err := e.DBrepo.GetUserInfo(userID)
	if err != nil {
		return "err", err
	}
	resData, err := e.DBrepo.GetRanking(userID, resUserInfo.UserCountry)
	if err != nil {
		return "err", err
	}
	// ここでresDataをJSONに変換する
	encodeJson,err := json.Marshal(resData)
	if err != nil {
		return "err", err
	}
	// ここでJSONを返す
	return string(encodeJson), nil
}