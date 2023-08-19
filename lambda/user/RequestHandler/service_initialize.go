package RequestHandler

import (
	"encoding/json"
	"log"
	"user/Data"
)

func (s APIFunctions) ServiceInitialize() (string, error) {
	if err := s.db_repo.AutoMigrate(); err != nil {
		return "Failed ServiceInitialize", err
	}
	// APIリクエスト制限を入れておく
	api_request_configs := []Data.ApiConfig{
		{
			AccountType:                 "Free",
			RefreshArticleInterval:      10,
			FetchArticleRequestInterval: 10,
			FetchArticleRequestLimit:    5000,
			FetchTrendRequestInterval:   5000,
			FetchTrendRequestLimit:      5000,
		},
		{
			AccountType:                 "Premium",
			RefreshArticleInterval:      1,
			FetchArticleRequestInterval: 1,
			FetchArticleRequestLimit:    100000,
			FetchTrendRequestInterval:   100000,
			FetchTrendRequestLimit:      100000,
		},
	}
	for _, api_request_config := range api_request_configs {
		api_request_config_json, err := json.Marshal(api_request_config)
		if err != nil {
			log.Println("Failed api_request_config Unmarshal error : ", err)
			return "Failed ServiceInitialize", err
		}
		res, err := s.ModifyAPIRequestLimit("Add", string(api_request_config_json))
		if err != nil || res != "Success ModifyAPIRequestLimit" {
			return "Failed ServiceInitialize", err
		}
	}
	// 後はExploreCategoryを入れておく
	explore_categories := []Data.ExploreCategory{
		{
			CategoryName: "Game",
		},
		{
			CategoryName: "News",
		},
		{
			CategoryName: "Entertainment",
		},
		{
			CategoryName: "IT",
		},
	}
	// DBRepoでカテゴリーを追加する
	for _, explore_category := range explore_categories {
		if err := s.db_repo.ModifyExploreCategory("Add",explore_category); err != nil {
			return "Failed ServiceInitialize", err
		}
	}
	return "Success ServiceInitialize", nil
}
