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
	dbRepo := DBRepo.DBRepoImpl{}
	if err := dbRepo.ConnectDB(true); err != nil {
		t.Fatal(err)
	}
	if err := dbRepo.AutoMigrate(); err != nil {
		t.Fatal(err)
	}

	// テスト用にDBにデータを入れておく
	site := Data.WebSite{
		SiteName:   "test",
		SiteURL:    "testSiteURL",
		SiteRssURL: "testRssURL",
	}
	dbSite := DBRepo.Site{
		SiteName: site.SiteName,
		SiteUrl:  site.SiteURL,
		RssUrl:   site.SiteRssURL,
	}
	DBRepo.DBMS.Create(&dbSite)
	article := Data.Article{
		Title: "test",
		Link:  "testArticleLink",
		Site:  site.SiteName,
	}
	if err := DBRepo.DBMS.Model(&dbSite).Association("SiteArticles").Append(&DBRepo.Article{
		Title: article.Title,
		Url:   article.Link,
	}); err != nil {
		t.Fatal(err)
	}

	// テスト用オブジェクトを用意する
	test_user_id := "test"
	now := time.Now()
	// RFC3339でフォーマットする
	nowStr := now.Format(time.RFC3339)
	readActJson, _ := json.Marshal(Data.ReadHistory{
		Link:     "test",
		AccessAt: nowStr,
	})
	userCfg := Data.UserConfig{
		UserName:     "test",
		UserUniqueID: test_user_id,
		ClientConfig: Data.ClientConfig{
			ApiConfig: Data.ApiConfig{
				RefreshArticleInterval: 10,
			},
		},
	}
	userConfigJson, _ := json.Marshal(userCfg)
	testWebSiteJson, _ := json.Marshal(site)
	testArticleJson, _ := json.Marshal(article)

	// 正解データを用意する
	configSyncResJson, _ := json.Marshal(Data.ConfigSyncResponse{
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
	ExpectedConfigSync := string(configSyncResJson)
	ExpectedRegisterUser := "Success RegisterUser"
	ExpectedReportReadActivity := "Success ReportReadActivity"
	ExpectedUpdateConfig := "Success UpdateConfig"
	ExpectedModifyFavoriteSite := "Success ModifyFavoriteSite"
	ExpectedModifyFavoriteArticle := "Success ModifyFavoriteArticle"
	apiRequestLimitCfgJson, _ := json.Marshal(Data.ApiConfig{
		RefreshArticleInterval: 10,
	})
	ExpectedApiRequestLimitCfg := string(apiRequestLimitCfgJson)
	searchHistoryJson, _ := json.Marshal([]string{"test"})
	ExpectedSearchHistoryJson := string(searchHistoryJson)
	UpdateApiConfig := Data.ApiConfig{
		RefreshArticleInterval: 20,
	}
	updateApiConfigJson, _ := json.Marshal(UpdateApiConfig)
	ExpectedUpdateApiConfigJson := "Success UpdateAPIRequestLimit"
	DeletedUserDataIsScoped, _ := json.Marshal(true)
	ExpectedDeleteUserData := "Success DeleteUserData"

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
				repo: dbRepo,
				ip:   "",
			},
			args: args{
				requestType:    "RegisterUser",
				userId:         test_user_id,
				argumentJson_1: string(userConfigJson),
				argumentJson_2: "",
			},
			want:    string(ExpectedRegisterUser),
			wantErr: false,
		},
		// GenUserIDは単純なのでテストしない
		{
			name: "ConfigSync",
			fields: fields{
				repo: dbRepo,
				ip:   "",
			},
			args: args{
				requestType:    "ConfigSync",
				userId:         test_user_id,
				argumentJson_1: "",
				argumentJson_2: "",
			},
			want:    ExpectedConfigSync,
			wantErr: false,
		},
		{
			name: "ReportReadActivity",
			fields: fields{
				repo: dbRepo,
				ip:   "",
			},
			args: args{
				requestType:    "ReportReadActivity",
				userId:         test_user_id,
				argumentJson_1: string(readActJson),
				argumentJson_2: "",
			},
			want:    string(ExpectedReportReadActivity),
			wantErr: false,
		},
		{
			name: "UpdateConfig",
			fields: fields{
				repo: dbRepo,
				ip:   "",
			},
			args: args{
				requestType:    "UpdateConfig",
				userId:         test_user_id,
				argumentJson_1: string(userConfigJson),
				argumentJson_2: "",
			},
			want:    string(ExpectedUpdateConfig),
			wantErr: false,
		},
		{
			name: "ModifySearchHistory",
			fields: fields{
				repo: dbRepo,
				ip:   "",
			},
			args: args{
				requestType:    "ModifySearchHistory",
				userId:         test_user_id,
				argumentJson_1: "test",
				argumentJson_2: "true",
			},
			want:    ExpectedSearchHistoryJson,
			wantErr: false,
		},
		{
			name: "ModifyFavoriteSite",
			fields: fields{
				repo: dbRepo,
				ip:   "",
			},
			args: args{
				requestType:    "ModifyFavoriteSite",
				userId:         test_user_id,
				argumentJson_1: string(testWebSiteJson),
				argumentJson_2: "true",
			},
			want:    string(ExpectedModifyFavoriteSite),
			wantErr: false,
		},
		{
			name: "ModifyFavoriteArticle",
			fields: fields{
				repo: dbRepo,
				ip:   "",
			},
			args: args{
				requestType:    "ModifyFavoriteArticle",
				userId:         test_user_id,
				argumentJson_1: string(testArticleJson),
				argumentJson_2: "true",
			},
			want:    string(ExpectedModifyFavoriteArticle),
			wantErr: false,
		},
		{
			name: "GetApiRequestLimit",
			fields: fields{
				repo: dbRepo,
				ip:   "",
			},
			args: args{
				requestType:    "GetAPIRequestLimit",
				userId:         test_user_id,
				argumentJson_1: "",
				argumentJson_2: "",
			},
			want:    ExpectedApiRequestLimitCfg,
			wantErr: false,
		},
		// UpdateApiRequestLimit
		{
			name: "UpdateApiRequestLimit",
			fields: fields{
				repo: dbRepo,
				ip:   "",
			},
			args: args{
				requestType:    "UpdateAPIRequestLimit",
				userId:         test_user_id,
				argumentJson_1: string(updateApiConfigJson),
				argumentJson_2: "",
			},
			want:    ExpectedUpdateApiConfigJson,
			wantErr: false,
		},
		// DeleteUserData
		{
			name: "DeleteUserData",
			fields: fields{
				repo: dbRepo,
				ip:   "",
			},
			args: args{
				requestType:    "DeleteUserData",
				userId:         test_user_id,
				argumentJson_1: string(DeletedUserDataIsScoped),
				argumentJson_2: "",
			},
			want:    ExpectedDeleteUserData,
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
