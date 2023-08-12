package APIFunction

import "encoding/json"

func (functions *APIFunctions) DeleteSite(access_ip string, site_url_json string, is_unscoped_json string) (string, error) {
	// site_url_jsonをjsonから変換する
	var site_url string
	if err := json.Unmarshal([]byte(site_url_json), &site_url); err != nil {
		return "", err
	}
	// is_unscoped_jsonをjsonから変換する
	var is_unscoped bool
	if err := json.Unmarshal([]byte(is_unscoped_json), &is_unscoped); err != nil {
		return "", err
	}
	// サイトを削除する
	if is_unscoped {
		if err := functions.DBRepo.DeleteSiteByUnscoped(site_url); err != nil {
			return "", err
		}
		return "Success DeleteSite", nil
	} else {
		if err := functions.DBRepo.DeleteSite(site_url); err != nil {
			return "", err
		}
		return "Success DeleteSite", nil
	}
}
