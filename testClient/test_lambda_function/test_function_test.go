package test_lambda_function

import (
	// "encoding/json"
	"testing"
)

func TestApiUser(t *testing.T) {
	t.Run("TestApiUserPart1", func(t *testing.T) {
		result,err := LambdaApiTest()
		if err != nil || !result  {
			t.Errorf("Error: %s", err.Error())
		}
	})
}
