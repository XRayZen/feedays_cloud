package RequestHandler

import (
	"encoding/json"
	"testing"
	"time"
	"user/DBRepo"
	"user/Data"
)

// 正常系のテスト
func TestNormalRequestHandler(t *testing.T) {
	db_repo := DBRepo.DBRepoImpl{}
	if err := db_repo.ConnectDB(true); err != nil {
		t.Fatal(err)
	}
	if err := db_repo.AutoMigrate(); err != nil {
		t.Fatal(err)
	}

	// テスト用にDBにデータを入れておく
	site := Data.WebSite{
		SiteName:   "test",
		SiteURL:    "testSiteURL",
		SiteRssURL: "testRssURL",
	}
	db_site := DBRepo.Site{
		SiteName: site.SiteName,
		SiteUrl:  site.SiteURL,
		RssUrl:   site.SiteRssURL,
	}
	DBRepo.DBMS.Create(&db_site)
	article := Data.Article{
		Title: "test",
		Link:  "testArticleLink",
		Site:  site.SiteName,
	}
	if err := DBRepo.DBMS.Model(&db_site).Association("SiteArticles").Append(&DBRepo.Article{
		Title: article.Title,
		Url:   article.Link,
	}); err != nil {
		t.Fatal(err)
	}

	// テスト用オブジェクトを用意する
	test_user_id := "test"
	now := time.Now()
	// RFC3339でフォーマットする
	now_str := now.Format(time.RFC3339)
	read_act_json, _ := json.Marshal(Data.ReadHistory{
		Link:     "test",
		AccessAt: now_str,
	})
	user_config := Data.UserConfig{
		UserName:     "test",
		UserUniqueID: test_user_id,
		ClientConfig: Data.ClientConfig{
			ApiConfig: Data.ApiConfig{
				RefreshArticleInterval: 10,
			},
		},
	}
	user_config_json, _ := json.Marshal(user_config)
	test_web_site_json, _ := json.Marshal(site)
	test_article_json, _ := json.Marshal(article)

	// 正解データを用意する
	config_sync_response_json, _ := json.Marshal(Data.ConfigSyncResponse{
		ResponseType: "accept",
		UserConfig: Data.UserConfig{
			UserName:     "test",
			UserUniqueID: test_user_id,
			ClientConfig: Data.ClientConfig{
				ApiConfig: Data.ApiConfig{
					RefreshArticleInterval: 10,
				},
			},
		},
		Error: "",
	})
	expected_config_sync := string(config_sync_response_json)
	expected_register_user := "Success RegisterUser"
	expected_report_read_activity := "Success ReportReadActivity"
	expected_update_config := "Success UpdateConfig"
	expected_modify_favorite_site := "Success ModifyFavoriteSite"
	expected_modify_favorite_article := "Success ModifyFavoriteArticle"
	api_request_limit_config_json, _ := json.Marshal(Data.ApiConfig{
		RefreshArticleInterval: 10,
	})
	expected_api_request_limit_config := string(api_request_limit_config_json)
	search_history_json, _ := json.Marshal([]string{"test"})
	expected_search_history_json := string(search_history_json)
	update_api_config := Data.ApiConfig{
		RefreshArticleInterval: 20,
	}
	update_api_config_json, _ := json.Marshal(update_api_config)
	expected_update_api_config_json := "Success UpdateAPIRequestLimit"
	deleted_user_data_is_scoped, _ := json.Marshal(true)
	expected_delete_user_data := "Success DeleteUserData"

	// テスト引数
	type fields struct {
		repo DBRepo.DBRepo
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
		wantErr  bool // エラーが発生するかどうか
	}{
		{
			name: "RegisterUser",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "RegisterUser",
				userId:         test_user_id,
				argumentJson_1: string(user_config_json),
				argumentJson_2: "",
			},
			want:    string(expected_register_user),
			wantErr: false,
		},
		// GenUserIDは単純なのでテストしない
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
			want:    expected_config_sync,
			wantErr: false,
		},
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
			want:    string(expected_report_read_activity),
			wantErr: false,
		},
		{
			name: "UpdateConfig",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "UpdateConfig",
				userId:         test_user_id,
				argumentJson_1: string(user_config_json),
				argumentJson_2: "",
			},
			want:    string(expected_update_config),
			wantErr: false,
		},
		{
			name: "ModifySearchHistory",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "ModifySearchHistory",
				userId:         test_user_id,
				argumentJson_1: "test",
				argumentJson_2: "true",
			},
			want:    expected_search_history_json,
			wantErr: false,
		},
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
			want:    string(expected_modify_favorite_site),
			wantErr: false,
		},
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
			want:    string(expected_modify_favorite_article),
			wantErr: false,
		},
		{
			name: "GetApiRequestLimit",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "GetAPIRequestLimit",
				userId:         test_user_id,
				argumentJson_1: "",
				argumentJson_2: "",
			},
			want:    expected_api_request_limit_config,
			wantErr: false,
		},
		// UpdateApiRequestLimit
		{
			name: "UpdateApiRequestLimit",
			fields: fields{
				repo: db_repo,
				ip:   "",
			},
			args: args{
				requestType:    "UpdateAPIRequestLimit",
				userId:         test_user_id,
				argumentJson_1: string(update_api_config_json),
				argumentJson_2: "",
			},
			want:    expected_update_api_config_json,
			wantErr: false,
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
			want:    expected_delete_user_data,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRequestType(tt.fields.ip, tt.fields.repo,
				tt.args.requestType, tt.args.userId, tt.args.argumentJson_1, tt.args.argumentJson_2)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestHandler.HandleRequest() errorType: %v error = %v, wantErr %v", tt.args.requestType, err, tt.wantErr)
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
			if got != tt.want {
				t.Errorf("RequestHandler.HandleRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
