package APIFunction

import (
	"encoding/json"
	"errors"
	"read/Data"
)

func (s APIFunctions) Search(access_ip string, user_id string, request_argument_json1 string) (string, error) {
	var result = Data.SearchResult{}
	result.ApiResponse = "" // 警告されるから一旦入れておく
	// リクエストをjsonから変換する
	var apiSearchRequest = Data.ApiSearchRequest{}
	if err := json.Unmarshal([]byte(request_argument_json1), &apiSearchRequest); err != nil {
		return "", err
	}
	// Wordが空文字の場合はエラー
	if apiSearchRequest.Word == "" {
		return "", errors.New("Word is empty")
	}
	// ワードがURLの場合
	if apiSearchRequest.Word[:4] == "http" {
		// サイトURLをキーにDBに該当するサイトがあるか確認する
		if s.DBRepo.IsExistSite(apiSearchRequest.Word) {
			// サイトURLをキーにDBに該当するサイトを返す
			webSite, err := s.DBRepo.FetchSite(apiSearchRequest.Word)
			if err != nil {
				return "", err
			}
			// 検索結果を返す
			result = Data.SearchResult{
				ApiResponse:     "accept",
				ResponseMessage: "success",
				ResultType:      "found",
				SearchType:      apiSearchRequest.SearchType,
				Websites: []Data.WebSite{
					webSite,
				},
			}
		} else {
			// 存在しない場合は新規サイト処理
			site, articles, err := newSite(apiSearchRequest.Word)
			if err != nil {
				return "", err
			}
			// 新規サイトをDBに登録する
			err = s.DBRepo.RegisterSite(site, articles)
			if err != nil {
				return "", err
			}
			// 検索結果を返す
			result = Data.SearchResult{
				ApiResponse:     "accept",
				ResponseMessage: "success",
				ResultType:      "new site",
				SearchType:      apiSearchRequest.SearchType,
				Websites: []Data.WebSite{
					site,
				},
			}
		}
	} else {
		// キーワード検索を行う
		articles, err := s.DBRepo.SearchArticlesByKeyword(apiSearchRequest.Word)
		if err != nil {
			return "", err
		}
		// 記事がゼロ件の場合はResultTypeをnoneにする
		if len(articles) == 0 {
			result = Data.SearchResult{
				ApiResponse:     "accept",
				ResponseMessage: "success",
				ResultType:      "none",
				SearchType:      apiSearchRequest.SearchType,
			}
		} else {
			// 検索結果を返す
			result = Data.SearchResult{
				ApiResponse:     "accept",
				ResponseMessage: "success",
				ResultType:      "found",
				SearchType:      apiSearchRequest.SearchType,
				Articles:        articles,
			}
		}
	}
	// webSiteをjsonに変換する
	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(resultJson), nil
}
