package test_lambda_function

import (
	// "encoding/json"
	"testing"
)

func TestApi(t *testing.T) {
	t.Run("TestApi", func(t *testing.T) {
		result,err := LambdaApiTest()
		if err != nil || !result  {
			t.Errorf("Error: %s", err.Error())
		}
	})
}
