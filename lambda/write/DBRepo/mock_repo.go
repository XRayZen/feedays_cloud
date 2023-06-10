package DBRepo

import "write/Data"


type MockDBRepo struct {
}

func (repo MockDBRepo) GetUserInfo(userId string) (Data.WebSite, error) {
	return Data.WebSite{}, nil
}


