package RequestHandler

// ハンドラーからゆーざーID とリクエストタイプを受け取って、DBからデータを取得して返す
import (
	"encoding/json"
	"read/Repo"
)


func GetExploreCategories(DBrepo Repo.DBRepository,userID string) (ExploreJson string, err error) {
	//DBからデータを取得する
	res_user_info, err := DBrepo.SearchUserConfig(userID,false)
	if err != nil {
		return "err", err
	}
	res_data, err := DBrepo.FetchExploreCategories(res_user_info.Country)
	if err != nil {
		return "err", err
	}
	// ここでresDataをJSONに変換する
	encode_json,err := json.Marshal(res_data)
	if err != nil {
		return "err", err
	}
	// ここでJSONを返す
	return string(encode_json), nil
}






