package RequestHandler

import (
	"encoding/json"
	"user/Data"
	"testing"
	"time"
	"user/DBRepo"
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
		Link: "test",
		AccessAt: nowStr,
	})
	userConfigJson, _ := json.Marshal(Data.UserConfig{
		UserName:     "test",
		UserUniqueID: test_user_id,
	})
	testWebSiteJson, _ := json.Marshal(site)
	testArticleJson, _ := json.Marshal(article)
	// 正解を用意する
	configSyncResJson, _ := json.Marshal(Data.ConfigSyncResponse{
		UserConfig: Data.UserConfig{
			UserName:     "test",
			UserUniqueID: test_user_id,
		},
		ResponseType: "accept",
		Error:        "",
	})
	apiResRegisterUser, _ := GenAPIResponse("accept", "Success RegisterUser", "")
	apiResReportReadActivity, _ := GenAPIResponse("accept", "Success ReportReadActivity", "")
	apiResUpdateConfig, _ := GenAPIResponse("accept", "Success UpdateConfig", "")
	apiResModifyFavoriteSite, _ := GenAPIResponse("accept", "Success ModifyFavoriteSite", "")
	apiResModifyFavoriteArticle, _ := GenAPIResponse("accept", "Success ModifyFavoriteArticle", "")
	apiRequestLimitCfg, _ := json.Marshal(Data.ApiConfig{})
	apiRequestLimitCfgJson := string(apiRequestLimitCfg)
	searchHistoryJson, _ := json.Marshal([]string{"test"})
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
			want:    string(apiResRegisterUser),
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
			want:    string(configSyncResJson),
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
			want:    string(apiResReportReadActivity),
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
			want:    string(apiResUpdateConfig),
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
			want:    string(searchHistoryJson),
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
			want:    string(apiResModifyFavoriteSite),
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
			want:    string(apiResModifyFavoriteArticle),
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
				userId:         "test",
				argumentJson_1: "",
				argumentJson_2: "",
			},
			want:    apiRequestLimitCfgJson,
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
			if got != tt.want {
				t.Errorf("RequestHandler.HandleRequest() = %v, want %v", got, tt.want)
			}
		})
	}

}
