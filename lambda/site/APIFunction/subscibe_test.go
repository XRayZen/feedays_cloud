package APIFunction

import (
	"encoding/json"
	"site/Data"
	"site/Repo"
	"testing"
)

func TestSubscribeSite(t *testing.T) {
	// サイトが登録されていたら購読を登録する
	// サイトが登録されていなかったら登録処理をしてから購読を登録する
	// Feedの時間はRFC3339形式で返す
	type args struct {
		access_ip   string
		user_id     string
		webSite     Data.WebSite
		isSubscribe bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// サイトが登録されていたら購読を登録する
		{
			name: "サイトが登録されていたら購読を登録する",
			args: args{
				access_ip: "",
				user_id:     "test_user_id",
				webSite: Data.WebSite{
					SiteURL:      "https://automaton-media.com/",
					SiteRssURL:   "https://automaton-media.com/feed/",
					SiteName:    "AUTOMATON",
				},
				isSubscribe: true,
			},
			want:    "Success Subscribe Site",
			wantErr: false,
		},
		// サイトが登録されていなかったら登録処理をしてから購読を登録する
		{
			name: "サイトが登録されていなかったら登録処理をしてから購読を登録する",
			args: args{
				access_ip: "",
				user_id:     "test_user_id",
				webSite: Data.WebSite{
					SiteURL:      "https://www.4gamer.net/",
					SiteRssURL:   "https://www.4gamer.net/rss/index.xml",
					SiteName:    "4Gamer.net",
				},
				isSubscribe: true,
			},
			want:    "Success Register Site",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		// テストデータをJSONに変換する
		request_argument_json1, err := json.Marshal(tt.args.webSite)
		if err != nil {
			t.Errorf("SubscribeSite() json.Marshal() error = %v", err)
			return
		}
		request_argument_json2, err := json.Marshal(tt.args.isSubscribe)
		if err != nil {
			t.Errorf("SubscribeSite() json.Marshal() error = %v", err)
			return
		}
		t.Run(tt.name, func(t *testing.T) {
			s := APIFunctions{
				DBRepo: Repo.MockDBRepo{},
			}
			got, err := s.SubscribeSite(tt.args.access_ip, tt.args.user_id, string(request_argument_json1), string(request_argument_json2))
			if (err != nil) != tt.wantErr {
				t.Errorf("SubscribeSite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				// テストデータをJSONに変換する
				request_argument_json1, err := json.Marshal(tt.args.webSite)
				if err != nil {
					t.Errorf("SubscribeSite() json.Marshal() error = %v", err)
					return
				}
				request_argument_json2, err := json.Marshal(tt.args.isSubscribe)
				if err != nil {
					t.Errorf("SubscribeSite() json.Marshal() error = %v", err)
					return
				}
				t.Errorf("SubscribeSite() = %v, want %v", got, tt.want)
				t.Errorf("SubscribeSite() request_argument_json1 = %v", string(request_argument_json1))
				t.Errorf("SubscribeSite() request_argument_json2 = %v", string(request_argument_json2))
			}
		})
	}
}
