package helpers

import (
	"github.com/dghubble/sling"
	"github.com/kpango/glg"
)

func Get(path string) map[string]interface{} {
	temp := make(map[string]interface{})
	_, err := sling.New().Get(path).Receive(temp, temp)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}

func Post(path string, body interface{}, token string) map[string]interface{} {
	temp := make(map[string]interface{})
	var err error
	if token == "" {
		_, err = sling.New().Post(path).BodyForm(body).Receive(temp, temp)
	} else {
		_, err = sling.New().Post(path).BodyForm(body).Set("x-gizo-token", token).Receive(temp, temp)
	}
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}

func Patch(path, token string) map[string]interface{} {
	temp := make(map[string]interface{})
	_, err := sling.New().Patch(path).Set("x-gizo-token", token).Receive(temp, temp)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}

func Delete(path, token string) map[string]interface{} {
	temp := make(map[string]interface{})
	_, err := sling.New().Delete(path).Set("x-gizo-token", token).Receive(temp, temp)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}
