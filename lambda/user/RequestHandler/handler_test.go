package RequestHandler

import (
	"encoding/json"
	"testing"
	"time"
	"user/DbRepo"
	"user/Data"
)

// リクエストタイプの正常系・異常系テスト
func TestNormalRequestHandler(t *testing.T) {
	db_repo := DbRepo.DBRepoImpl{}
	// DBにモックモードで接続する（SQLite：メモリ上にDBを作成する）
	db_repo.ConnectDB(true)
	db_repo.AutoMigrate()

	// テスト用にDBにデータを入れておく
	site := Data.WebSite{
		SiteName:   "test",
		SiteURL:    "testSiteURL",
		SiteRssURL: "testRssURL",
	}
	db_site := DbRepo.Site{
		SiteName: site.SiteName,
		SiteUrl:  site.SiteURL,
		RssUrl:   site.SiteRssURL,
	}
	DbRepo.DBMS.Create(&db_site)
	article := Data.Article{
		Title: "test",
		Link:  "testArticleLink",
		Site:  site.SiteName,
	}
	DbRepo.DBMS.Model(&db_site).Association("SiteArticles").Append(&DbRepo.Article{
		Title: article.Title,
		Url:   article.Link,
	})

	// テスト用オブジェクトを用意する
	test_user_id := "test"
	now := time.Now()
	// RFC3339でフォーマットする
	now_str := now.Format(time.RFC3339)
	read_act_json, _ := json.Marshal(Data.ReadHistory{
		Link:     "test",
		AccessAt: now_str,
	})
	register_user_config := Data.UserConfig{
		UserName:     "test",
		UserUniqueID: test_user_id,
		AccountType:  "Free",
		ClientConfig: Data.ClientConfig{
			UiConfig: Data.UiConfig{
				ThemeMode: "light",
			},
		},
	}
	register_user_config_json, _ := json.Marshal(register_user_config)
	test_web_site_json, _ := json.Marshal(site)
	test_article_json, _ := json.Marshal(article)

	// 期待する出力データを定義
	updateUiConfigJson, _ := json.Marshal(Data.UserConfig{
		UserName:     "test",
		UserUniqueID: test_user_id,
		AccountType:  "Free",
		ClientConfig: Data.ClientConfig{
			UiConfig: Data.UiConfig{
				ThemeMode: "dark",
			},
		},
	})
	expected_config_sync_response_json, _ := json.Marshal(Data.ConfigSyncResponse{
		ResponseType: "accept",
		UserConfig: register_user_config,
		Error: "",
	})
	expected_config_sync := string(expected_config_sync_response_json)
	// 入力と期待する出力
	search_history_json, _ := json.Marshal(Data.SearchHistory{
		SearchWord: "test",
		SearchAt:   now_str,
	})
	expected_search_history_json,_ := json.Marshal([]string{"test"})
	deleted_user_data_is_scoped, _ := json.Marshal(true)

	// テスト引数
	type fields struct {
		repo DbRepo.DBRepo
		ip   string
	}
	type args struct {
		requestType    string
		userId         string
		argumentJson_1 string
		argumentJson_2 string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		response string
		want     string
	}{
		// RegisterUser
		{
			name: "RegisterUser",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "RegisterUser",
				userId:         test_user_id,
				argumentJson_1: string(register_user_config_json),
				argumentJson_2: "",
			},
			want: "Success RegisterUser",
		},
		// GenUserIDは単純なのでテストしない
		// ConfigSync
		{
			name: "ConfigSync",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "ConfigSync",
				userId:         test_user_id,
				argumentJson_1: "",
				argumentJson_2: "",
			},
			want: expected_config_sync,
		},
		// ReportReadActivity
		{
			name: "ReportReadActivity",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "ReportReadActivity",
				userId:         test_user_id,
				argumentJson_1: string(read_act_json),
				argumentJson_2: "",
			},
			want: "Success ReportReadActivity",
		},
		// UpdateUiConfig
		{
			name: "UpdateUiConfig",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "UpdateUiConfig",
				userId:         test_user_id,
				argumentJson_1: string(updateUiConfigJson),
				argumentJson_2: "",
			},
			want: "Success UpdateUiConfig",
		},
		// ModifySearchHistory
		{
			name: "ModifySearchHistory",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "ModifySearchHistory",
				userId:         test_user_id,
				argumentJson_1: string(search_history_json),
				argumentJson_2: "true",
			},
			want: string(expected_search_history_json),
		},
		// ModifyFavoriteSite
		{
			name: "ModifyFavoriteSite",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "ModifyFavoriteSite",
				userId:         test_user_id,
				argumentJson_1: string(test_web_site_json),
				argumentJson_2: "true",
			},
			want: "Success ModifyFavoriteSite",
		},
		// ModifyFavoriteArticle
		{
			name: "ModifyFavoriteArticle",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "ModifyFavoriteArticle",
				userId:         test_user_id,
				argumentJson_1: string(test_article_json),
				argumentJson_2: "true",
			},
			want: "Success ModifyFavoriteArticle",
		},
		// DeleteUserData
		{
			name: "DeleteUserData",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "DeleteUserData",
				userId:         test_user_id,
				argumentJson_1: string(deleted_user_data_is_scoped),
				argumentJson_2: "",
			},
			want: "Success DeleteUserData",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRequestType(tt.fields.ip, tt.fields.repo,
				tt.args.requestType, tt.args.userId, tt.args.argumentJson_1, tt.args.argumentJson_2)
			if err != nil {
				t.Errorf("RequestHandler.HandleRequest() errorType: %v error = %v", tt.args.requestType, err)
				return
			}
			// 念の為にユーザーデータ削除の場合で、削除後のデータが空であることを確認する必要がある
			if tt.args.requestType == "DeleteUserData" {
				_, err = ParseRequestType(tt.fields.ip, tt.fields.repo,
					"ConfigSync", tt.args.userId, "", "")
				// エラーならテスト成功
				if err.Error() == "record not found" {
					return
				} else {
					t.Errorf("DeleteUserData failed")
					return
				}
			}
			// UpdateUiConfigの場合は、ConfigSyncを呼び出して更新されているか確認する
			if tt.args.requestType == "UpdateUiConfig" {
				// ConfigSyncを呼び出す
				config_sync_response_json, err := ParseRequestType(tt.fields.ip, tt.fields.repo,
					"ConfigSync", tt.args.userId, "", "")
				if err != nil {
					t.Errorf("UpdateUiConfig failed")
					return
				}
				config_sync := Data.ConfigSyncResponse{}
				if err := json.Unmarshal([]byte(config_sync_response_json), &config_sync); err != nil {
					t.Errorf("UpdateUiConfig failed")
					return
				}
				if config_sync.UserConfig.ClientConfig.UiConfig.ThemeMode != "dark" {
					t.Errorf("UpdateUiConfig failed")
					return
				}
			}
			if got != tt.want {
				t.Errorf("RequestHandler.HandleRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

// APIリクエスト制限の正常系・異常系テスト
func TestApiRequestLimit(t *testing.T) {
	// 最初にユーザーを登録する
	db_repo := DbRepo.DBRepoImpl{}
	// DBにモックモードで接続する（SQLite：メモリ上にDBを作成する）
	db_repo.ConnectDB(true)
	db_repo.AutoMigrate()
	db_repo.RegisterUser(Data.UserConfig{
		UserName:     "test",
		UserUniqueID: "test",
		AccountType:  "Free",
		Country:      "JP",
	})
	// テスト用オブジェクトを用意する
	input_add_api_config_json, _ := json.Marshal(Data.ApiConfig{
		AccountType:            "Free",
		RefreshArticleInterval: 10,
	})
	input_update_api_config_json, _ := json.Marshal(Data.ApiConfig{
		AccountType:            "Free",
		RefreshArticleInterval: 20,
	})
	// テストケース
	type args struct {
		requestType    string
		userId         string
		argumentJson_1 string
		argumentJson_2 string
	}
	tests := []struct {
		name        string
		args        args
		want        string
		is_want_err bool // エラーが発生するかどうか
	}{
		// AddApiRequestLimit
		{
			name: "AddApiRequestLimit",
			args: args{
				requestType:    "ModifyAPIRequestLimit",
				userId:         "test",
				argumentJson_1: "Add",
				argumentJson_2: string(input_add_api_config_json),
			},
			want:        "Success ModifyAPIRequestLimit",
			is_want_err: false,
		},
		// GetApiRequestLimit
		{
			name: "GetApiRequestLimit",
			args: args{
				requestType:    "GetAPIRequestLimit",
				userId:         "test",
				argumentJson_1: "",
				argumentJson_2: "",
			},
			want:        string(input_add_api_config_json),
			is_want_err: false,
		},
		// UpdateApiRequestLimit
		{
			name: "UpdateApiRequestLimit",
			args: args{
				requestType:    "ModifyAPIRequestLimit",
				userId:         "test",
				argumentJson_1: "Update",
				argumentJson_2: string(input_update_api_config_json),
			},
			want:        "Success ModifyAPIRequestLimit",
			is_want_err: false,
		},
		// DeleteApiRequestLimit
		{
			name: "DeleteApiRequestLimit",
			args: args{
				requestType:    "ModifyAPIRequestLimit",
				userId:         "test",
				argumentJson_1: "UnscopedDelete",
				argumentJson_2: string(input_add_api_config_json),
			},
			want:        "Success ModifyAPIRequestLimit",
			is_want_err: false,
		},
		// 異常系テスト
		// 削除されているのに読み込もうとする
		{
			name: "Failed GetApiRequestLimit",
			args: args{
				requestType:    "GetAPIRequestLimit",
				userId:         "test",
				argumentJson_1: "",
				argumentJson_2: "",
			},
			want:        "record not found",
			is_want_err: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRequestType("", db_repo,
				tt.args.requestType, tt.args.userId, tt.args.argumentJson_1, tt.args.argumentJson_2)
			if err != nil && tt.is_want_err {
				if err.Error() == tt.want {
					return // 期待したエラーが発生しているのでテスト成功
				}
				// 期待したエラーが発生していないのでテスト失敗
				t.Errorf("RequestHandler.HandleRequest() errorType: %v error = %v, wantErr %v", tt.args.requestType, err, tt.is_want_err)
			} else if err != nil && !tt.is_want_err {
				t.Errorf("RequestHandler.HandleRequest() errorType: %v error = %v, wantErr %v", tt.args.requestType, err, tt.is_want_err)
			}
			if got != tt.want {
				t.Errorf(err.Error())
				t.Errorf("RequestHandler.HandleRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

// サービスの初期化と終了のテスト
func TestServiceInitAndClose(t *testing.T) {
	db_repo := DbRepo.DBRepoImpl{}
	// DBにモックモードで接続する（SQLite：メモリ上にDBを作成する）
	db_repo.ConnectDB(true)
	// テストケース
	tests := []struct {
		name          string
		request_type  string // テストするリクエストタイプ
		want          string
		is_expect_err bool   // エラーが発生するかどうか
		expect_err    string // 期待するエラー内容
	}{
		// ServiceInitialize
		{
			name:          "ServiceInitialize",
			request_type:  "ServiceInitialize",
			want:          "Success ServiceInitialize",
			is_expect_err: false,
			expect_err:    "",
		},
		// ServiceFinalize
		{
			name:          "ServiceFinalize",
			request_type:  "ServiceFinalize",
			want:          "Success ServiceFinalize",
			is_expect_err: false,
			expect_err:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRequestType("", db_repo,
				tt.request_type, "", "", "")
			if err != nil && !tt.is_expect_err {
				t.Errorf("RequestHandler.HandleRequest() errorType: %v error = %v, wantErr %v", tt.request_type, err, tt.is_expect_err)
			}
			if got != tt.want {
				t.Errorf("RequestHandler.HandleRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
