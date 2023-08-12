package APIFunction

import "encoding/json"

// サイトにExploreCategoryを追加する為にサイトのカテゴリを変更する
func (functions *APIFunctions) ChangeSiteCategory(access_ip string, user_id string, site_url_json string, category_name_json string) (string, error) {
	// site_url_jsonをjsonから変換する
	var site_url string
	if err := json.Unmarshal([]byte(site_url_json), &site_url); err != nil {
		return "", err
	}
	// category_name_jsonをjsonから変換する
	var category_name string
	if err := json.Unmarshal([]byte(category_name_json), &category_name); err != nil {
		return "", err
	}
	// category_name_jsonをデコードする
	// サイトのカテゴリを変更する
	if err := functions.DBRepo.ChangeSiteCategory(user_id, site_url, category_name); err != nil {
		return "", err
	}
	return "Success ChangeSiteCategory", nil
}
