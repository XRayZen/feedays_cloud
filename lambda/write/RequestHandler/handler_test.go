package RequestHandler

import (
	"encoding/json"
	"testing"
	"write/DBRepo"
	"read/Data"
)

// 正常系のテスト
func TestNormalRequestHandler(t *testing.T) {
	// テスト用オブジェクトを用意する
	identInfoJson, _ := json.Marshal(Data.UserAccessIdentInfo{
		UUid: "test",
		AccessPlatform: "test",
		PlatformType: "test",
		Brand: "test",
		DeviceName: "test",
		OsVersion: "test",
		IsPhysics: false,
	})
	readActJson, _ := json.Marshal(Data.ReadActivity{
		UserID: "test",
		Title:  "test",
		Link:   "test",
		Type:   "test",
	})
	userConfigJson, _ := json.Marshal(Data.UserConfig{})
	testWebSiteJson, _ := json.Marshal(Data.WebSite{})
	testArticleJson, _ := json.Marshal(Data.Article{})
	// 正解を用意する
	configSyncResJson, _ := json.Marshal(Data.ConfigSyncResponse{
		UserConfig:   Data.UserConfig{},
		ResponseType: "accept",
		Error:        "",
	})
	apiResRegisterUser, _ := GenAPIResponse("accept", "Success RegisterUser", "")
	apiResReportReadActivity, _ := GenAPIResponse("accept", "Success ReportReadActivity", "")
	apiResUpdateConfig, _ := GenAPIResponse("accept", "Success UpdateConfig", "")
	apiResModifyFavoriteSite, _ := GenAPIResponse("accept", "Success ModifyFavoriteSite", "")
	apiResModifyFavoriteArticle, _ := GenAPIResponse("accept", "Success ModifyFavoriteArticle", "")
	apiRequestLimitCfg, _ := json.Marshal(Data.ApiRequestLimitConfig{})
	apiRequestLimitCfgJson := string(apiRequestLimitCfg)
	searchHistoryJson, _ := json.Marshal([]string{})

	// テスト引数
	type fields struct {
		repo DBRepo.MockDBRepo
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
		// GenUserIDは単純なのでテストしない
		{
			name: "ConfigSync",
			fields: fields{
				repo: DBRepo.MockDBRepo{},
				ip:   "",
			},
			args: args{
				requestType:    "ConfigSync",
				userId:         "",
				argumentJson_1: string(identInfoJson),
				argumentJson_2: "",
			},
			want:    string(configSyncResJson),
			wantErr: false,
		},
		{
			name: "RegisterUser",
			fields: fields{
				repo: DBRepo.MockDBRepo{},
				ip:   "",
			},
			args: args{
				requestType:    "RegisterUser",
				userId:         "",
				argumentJson_1: string(userConfigJson),
				argumentJson_2: string(identInfoJson),
			},
			want:    string(apiResRegisterUser),
			wantErr: false,
		},
		{
			name: "ReportReadActivity",
			fields: fields{
				repo: DBRepo.MockDBRepo{},
				ip:   "",
			},
			args: args{
				requestType:    "ReportReadActivity",
				userId:         "",
				argumentJson_1: string(readActJson),
				argumentJson_2: "",
			},
			want:    string(apiResReportReadActivity),
			wantErr: false,
		},
		{
			name: "UpdateConfig",
			fields: fields{
				repo: DBRepo.MockDBRepo{},
				ip:   "",
			},
			args: args{
				requestType:    "UpdateConfig",
				userId:         "",
				argumentJson_1: string(userConfigJson),
				argumentJson_2: string(identInfoJson),
			},
			want:    string(apiResUpdateConfig),
			wantErr: false,
		},
		{
			name: "ModifySearchHistory",
			fields: fields{
				repo: DBRepo.MockDBRepo{},
				ip:   "",
			},
			args: args{
				requestType:    "ModifySearchHistory",
				userId:         "",
				argumentJson_1: "test",
				argumentJson_2: "true",
			},
			want:    string(searchHistoryJson),
			wantErr: false,
		},
		{
			name: "ModifyFavoriteSite",
			fields: fields{
				repo: DBRepo.MockDBRepo{},
				ip:   "",
			},
			args: args{
				requestType:    "ModifyFavoriteSite",
				userId:         "",
				argumentJson_1: string(testWebSiteJson),
				argumentJson_2: "true",
			},
			want:    string(apiResModifyFavoriteSite),
			wantErr: false,
		},
		{
			name: "ModifyFavoriteArticle",
			fields: fields{
				repo: DBRepo.MockDBRepo{},
				ip:   "",
			},
			args: args{
				requestType:    "ModifyFavoriteArticle",
				userId:         "",
				argumentJson_1: string(testArticleJson),
				argumentJson_2: "true",
			},
			want:    string(apiResModifyFavoriteArticle),
			wantErr: false,
		},
		{
			name: "GetApiRequestLimit",
			fields: fields{
				repo: DBRepo.MockDBRepo{},
				ip:   "",
			},
			args: args{
				requestType:    "GetAPIRequestLimit",
				userId:         "",
				argumentJson_1: string(identInfoJson),
				argumentJson_2: "",
			},
			want:    apiRequestLimitCfgJson,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// handler := &RequestHandler{
			// 	repo: tt.fields.repo,
			// 	ip:   tt.fields.ip,
			// }
			got, err := ParseRequestType(tt.fields.ip, tt.fields.repo,
				tt.args.requestType, tt.args.userId, tt.args.argumentJson_1, tt.args.argumentJson_2)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestHandler.HandleRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RequestHandler.HandleRequest() = %v, want %v", got, tt.want)
			}
		})
	}

}
