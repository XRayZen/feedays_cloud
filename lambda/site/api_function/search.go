package APIFunction

import (
	"encoding/json"
	"errors"
	"site/Data"
	"site/Repo"
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
		return "", errors.New("search word is empty")
	}
	switch apiSearchRequest.SearchType {
	case "URL":
		// URL検索
		// DBに存在しない場合は新規サイト処理
		res, err := searchByURL(s.DBRepo, apiSearchRequest)
		if err != nil {
			return "", err
		}
		result = res
	case "Keyword":
		res, err := searchByKeyword(s.DBRepo, apiSearchRequest)
		if err != nil {
			return "", err
		}
		result = res
	case "SiteName":
		res, err := searchBySiteName(s.DBRepo, apiSearchRequest)
		if err != nil {
			return "", err
		}
		result = res
	case "Category":
		res, err := searchByCategory(s.DBRepo, apiSearchRequest)
		if err != nil {
			return "", err
		}
		result = res
	}
	// 検索をアクテビティとして報告するが、今は実装しない
	// resultをjsonに変換する
	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(resultJson), nil
}

func searchByCategory(dBRepository Repo.DBRepository, apiSearchRequest Data.ApiSearchRequest) (result Data.SearchResult, err error) {
	db_result, err := dBRepository.SearchSiteByCategory(apiSearchRequest.Word)
	if err != nil {
		return Data.SearchResult{}, err
	}
	return Data.SearchResult{
		ApiResponse:     "accept",
		ResponseMessage: "success",
		ResultType:      "found",
		SearchType:      apiSearchRequest.SearchType,
		Websites:        db_result,
	}, nil
}

// URL検索
func searchByURL(repo Repo.DBRepository, apiSearchRequest Data.ApiSearchRequest) (result Data.SearchResult, err error) {
	// サイトURLをキーにDBに該当するサイトがあるか確認する
	if repo.IsExistSite(apiSearchRequest.Word) {
		// サイトURLをキーにDBに該当するサイトを返す
		webSite, err := repo.SearchSiteByUrl(apiSearchRequest.Word)
		if err != nil {
			return Data.SearchResult{}, err
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
		// サイトテーブルに存在しない場合は新規サイト処理
		site, articles, err := newSite(apiSearchRequest.Word)
		if err != nil {
			// RSSの取得に失敗した場合はエラーで汎用APIエラー応答にする
			return Data.SearchResult{}, err
		}
		if err := repo.RegisterSite(site, articles); err != nil {
			return Data.SearchResult{}, err
		}
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
	return result, nil
}

// 記事へのキーワード検索
func searchByKeyword(repo Repo.DBRepository, apiSearchRequest Data.ApiSearchRequest) (result Data.SearchResult, err error) {
	// キーワード検索は本来ならフリーアカウントならできないようにするべきだが、
	// 今回はフリーアカウントでもできるようにする
	// キーワード検索を行う
	articles, err := repo.SearchArticlesByKeyword(apiSearchRequest.Word)
	if err != nil {
		return Data.SearchResult{}, err
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
	return result, nil
}

// サイト名を検索
func searchBySiteName(repo Repo.DBRepository, apiSearchRequest Data.ApiSearchRequest) (result Data.SearchResult, err error) {
	res, err := repo.SearchSiteByName(apiSearchRequest.Word)
	if err != nil {
		return Data.SearchResult{}, err
	}
	// 検索結果を返す
	if len(res) == 0 {
		result = Data.SearchResult{
			ApiResponse:     "accept",
			ResponseMessage: "success",
			ResultType:      "none",
			SearchType:      apiSearchRequest.SearchType,
		}
	} else {
		result = Data.SearchResult{
			ApiResponse:     "accept",
			ResponseMessage: "success",
			ResultType:      "found",
			SearchType:      apiSearchRequest.SearchType,
			Websites:        res,
		}
	}
	return result, nil
}
