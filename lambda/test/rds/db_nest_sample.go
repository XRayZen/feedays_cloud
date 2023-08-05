package RDS

import "log"

type Client struct {
	//  ApiKey Index
	ID       int    `json:"id"`
	UserName string `json:"userName"`
	//  Client login hashed password
	Password string `json:"password"`
	//  ApiKeys
	APIKeys []*APIKey `json:"apiKeys"`
}

type APIKey struct {
	//  ApiKey Index
	ID int `json:"id"`
	//  ApiKey Value
	Key string `json:"key"`
	//  ApiKey Client Relation
	ClientID int `json:"clientID"`
	//  ApiKey Client Info
	Client *Client `json:"client"`
}

func DbNestedSampleTest() (bool, error) {
	db, err := GormConnect()
	if err != nil {
		return false, err
	}
	// テーブル作成
	err = db.AutoMigrate(&Client{}, &APIKey{})
	if err != nil {
		log.Println("AutoMigrate Error:", err)
		return false, err
	}

	clientOne := Client{
		UserName: "Client One",
	}
	db.Create(&clientOne)

	apiKeyOne := APIKey{
		Key:    "one",
		Client: &clientOne,
	}
	apiKeyTwo := APIKey{
		Key:    "two",
		Client: &clientOne,
	}

	db.Create(&apiKeyOne)
	db.Create(&apiKeyTwo)

	// Fetch from DB
	fetchedClient := Client{}

	db.Debug().Preload("APIKeys").Find(&fetchedClient, clientOne.ID)
	log.Println("ApiKey length:", len(fetchedClient.APIKeys))

	db.Delete(&clientOne)
	db.Delete(&apiKeyOne)
	db.Delete(&apiKeyTwo)
	return true, nil
}
