package chef

import (
	"io/ioutil"
	"os"

	"git.rn/devops/zhash.git"
)

func readFile(fn string) (string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	return string(buf), err
}

func getErrorMessage(responce zhash.Hash) string {
	var errText string
	errorMessage, _ := responce.GetStringSlice("error")
	for _, e := range errorMessage {
		errText += e
	}
	if errText == "" {
		errText = responce.String()
	}
	return errText
}
