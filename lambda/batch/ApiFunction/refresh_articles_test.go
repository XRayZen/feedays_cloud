package ApiFunction

import (
	"log"
	"batch/Repo"
	"testing"
	"time"
)

func TestRefreshArticles(t *testing.T) {
	// DBRepo
	mock_repo := Repo.MockDBRepo{}
	// implRepo := Repo.DBRepoImpl{}
	// 考えられるテストケースを網羅する
	// 更新する記事がない場合
	// 更新する記事がある場合
	type args struct {
		repo                 Repo.DBRepository
		MockSiteLastModified int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "更新する記事がない場合",
			args: args{
				repo:                 mock_repo,
				MockSiteLastModified: 10,
			},
			want:    "No Update",
			wantErr: false,
		},
		{
			name: "更新する記事がある場合",
			args: args{
				repo:                 mock_repo,
				MockSiteLastModified: 30,
			},
			want:    "BATCH RefreshArticles SUCCESS!",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Repo.MockSiteLastModified = tt.args.MockSiteLastModified
			got, err := RefreshArticles(tt.args.repo,15)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RefreshArticles() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_refreshSiteArticles(t *testing.T) {
	// 考えられるテストケースを網羅する
	// サイトが存在していて更新期限を過ぎている場合は
	// 記事を更新してDBに渡してクライアント側更新日時より新しい記事を返す
	// サイトが存在していて更新期限を過ぎていない場合は
	// 更新せずにクライアント側更新日時より新しい記事を返す
	type args struct {
		siteUrl            string
		db_repo            Repo.DBRepository
		intervalMinutes    int
		clientLastModified time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "サイトが存在していて更新期限を過ぎている場合",
			args: args{
				siteUrl:            "https://automaton-media.com/",
				db_repo:            Repo.MockDBRepo{},
				intervalMinutes:    5,
				clientLastModified: time.Now().Add(-time.Hour * 160),
			},
			// このテストケースはサイトの記事を更新してDBに渡してクライアント側更新日時より新しい記事を返す
			want:    2,
			wantErr: false,
		},
		{
			name: "サイトが存在していて更新期限を過ぎていない場合",
			args: args{
				siteUrl:            "https://gigazine.net/",
				db_repo:            Repo.MockDBRepo{},
				intervalMinutes:    1,
				clientLastModified: time.Now().Add(-time.Hour * 10),
			},
			// このテストケースは更新せずにクライアント側更新日時より新しい記事を返す
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RefreshArticles(tt.args.db_repo,15)
			if (err != nil) != tt.wantErr {
				t.Errorf("refreshSiteArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// デバッグ用ログ
			log.Println("got: ", len(got))
			log.Println("want: ", tt.want)
			if len(got) < tt.want {
				t.Errorf("refreshSiteArticles() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_isUpdateExpired(t *testing.T) {
	// 更新期限を過ぎている場合はtrueを返す
	// 更新期限を過ぎていない場合はfalseを返す
	type args struct {
		lastModified    time.Time
		intervalMinutes int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "更新期限を過ぎている場合",
			args: args{
				// 更新日時は10分前
				lastModified: time.Now().Add(-time.Minute * 10),
				// 更新期限は5分
				intervalMinutes: 5,
			},
			want: true,
		},
		{
			name: "更新期限を過ぎていない場合",
			args: args{
				// 更新日時は10分前
				lastModified: time.Now().Add(-time.Minute * 10),
				// 更新期限は15分
				intervalMinutes: 15,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		if got := isUpdateExpired(tt.args.lastModified, tt.args.intervalMinutes); got != tt.want {
			t.Errorf("%q. isUpdateExpired() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
