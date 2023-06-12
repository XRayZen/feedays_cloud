package RequestHandler

import (
	"encoding/json"
	"write/DBRepo"
	"write/Data"
	"testing"
)

func TestAPIFunctions_GetUserInfo(t *testing.T) {
	type fields struct {
		repo DBRepo.MockDBRepo
		ip   string
	}
	type args struct {
		userId string
	}



}


