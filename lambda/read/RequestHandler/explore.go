package RequestHandler

// ハンドラーからゆーざーID とリクエストタイプを受け取って、DBからデータを取得して返す
import (
	"encoding/json"
	"read/Repo"
)


func GetExploreCategories(DBrepo Repo.DBRepository,userID string) (ExploreJson string, err error) {
	//DBからデータを取得する
	resUserInfo, err := DBrepo.SearchUserConfig(userID,false)
	if err != nil {
		return "err", err
	}
	resData, err := DBrepo.FetchExploreCategories(resUserInfo.Country)
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






