package RequestHandler

import (
	// "errors"
	"encoding/json"
	"read/DBRepo"
	"read/Data"

	"testing"
)

// GetExploreCategoriesをテストする
func TestGetExploreCategories(t *testing.T) {
	type args struct {
		userID string
	}
	// 欲しいデータを定義する
	want :=Data.ExploreCategories{
		CategoryName:        "CategoryName",
	}

	tests := []struct {
		name    string
		args    args
		want    Data.ExploreCategories
		wantErr bool
	}{

		{
			name: "正常系",
			args: args{
				userID: "userID",
			},
			want:    want,
			wantErr: false,
		},
		// BUG:この関数で異常系は必要ない（何かを計算するわけではない）だから、異常系のテストは不要
		{
			name: "異常系",
			args: args{
				userID: "userID",
			},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ここでモックにDIする
			s := Explore{
				// DBrepo: DBRepo.DBRepoImpl{},
				DBrepo: DBRepo.MockDBRepo{},
			}
			got, err :=s.GetExploreCategories(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExploreCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			wat ,_:= json.Marshal(tt.want)
			if got != string(wat) {
				t.Errorf("GetExploreCategories() got = %v, want %v", got, tt.want)
			}
		})
	}
}
