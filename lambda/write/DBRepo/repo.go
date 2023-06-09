package DBRepo

import (
	"errors"
	"read/RequestType"
)

// DBRepo はDBにアクセスするためのインターフェース
type DBRepo interface {
	GetUserInfo(userId string) (string, error)

}