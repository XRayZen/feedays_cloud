package RequestHandler

import (
	"encoding/json"
	"write/Data"
)

func GenAPIResponse(responseType string, msg string, errorMsg string) (string, error) {
	res := Data.APIResponse{
		ResponseType: responseType,
		Message:      msg,
		Error:        errorMsg,
	}
	res_str, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(res_str), nil
}
