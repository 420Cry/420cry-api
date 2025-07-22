package constants

import "fmt"

func ReturnConstant(api string) (string, error) {
	var constant string
	if api == "googleUserApi" {
		constant = "https://www.googleapis.com/oauth2/v3/userinfo"
	}

	if constant == "" {
		return "", fmt.Errorf("no constant found")
	}

	return constant, nil

}
