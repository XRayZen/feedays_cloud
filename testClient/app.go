package main

import (
	"fmt"
	"testClient/test_lambda_function"
)

// APIがちゃんと動くかのテスト用
// ちゃんと変換が確認されていたら、フォルダごと削除してもいい
func main() {
	fmt.Println("Begin Api Test")
	result,err :=test_lambda_function.LambdaApiTest()
	if err != nil || !result  {
		fmt.Println("Error: ", err.Error())
	}
	fmt.Println("End Api Test")
}
