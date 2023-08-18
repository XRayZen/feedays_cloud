package ApiFunction

import (
	"encoding/json"
	"testing"
	"site/Data"
	"site/Repo"
)

func Test_Search(t *testing.T) {
	type args struct {
		access_ip        string
		user_id          string
		apiSearchRequest Data.ApiSearchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "URL検索でサイトが見つからない場合",
			args: args{
				access_ip: "",
				user_id:   "test",
				apiSearchRequest: Data.ApiSearchRequest{
					Word:        "https://jp.ign.com/",
					SearchType: "URL",
				},
			},
			// 新規サイト処理をしてnew siteを返す
			want:    "new site",
			wantErr: false,
		},
		{
			name: "URL検索でサイトが見つかる場合",
			args: args{
				access_ip:        "",
				user_id:          "test",
				apiSearchRequest: Data.ApiSearchRequest{
					Word:        "https://gigazine.net/",
					SearchType: "URL",
				},
			},
			// 結果にサイトを含めてfoundを返す
			want:    "found",
			wantErr: false,
		},
		{
			name: "キーワード検索で記事が見つからない場合",
			args: args{
				access_ip:        "",
				user_id:          "test",
				apiSearchRequest: Data.ApiSearchRequest{
					Word: 	  "test",
					SearchType: "Keyword",
				},
			},
			want:    "none",
			wantErr: false,
		},
		{
			name: "キーワード検索で記事が見つかる場合",
			args: args{
				access_ip:        "",
				user_id:          "test",
				apiSearchRequest: Data.ApiSearchRequest{
					Word: 	  "Found",
					SearchType: "Keyword",
				},
			},
			want:    "found",
			wantErr: false,
		},
		{
			name: "サイト検索でDBに存在した場合",
			args: args{
				access_ip:        "",
				user_id:          "test",
				apiSearchRequest: Data.ApiSearchRequest{
					Word: 	  "Found",
					SearchType: "SiteName",
				},
			},
			want:    "found",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 検索リクエストをJsonに変換
			search_json, _ := json.Marshal(tt.args.apiSearchRequest)
			// ここでモックへの依存性を注入
			functions := APIFunctions{
				DBRepo: Repo.MockDBRepo{},
			}
			got, err := functions.Search(tt.args.access_ip, tt.args.user_id, string(search_json))
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 検索結果jsonを検索結果構造体に変換
			var result = Data.SearchResult{}
			json.Unmarshal([]byte(got), &result)
			// 検索結果構造体の中身を検証
			if result.ResultType != tt.want {
				t.Errorf("Search() got = %v, want %v", got, tt.want)
			}
		})
	}

}
