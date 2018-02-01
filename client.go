package chef

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/zazab/zhash"
)

type Client struct {
	Name      string `json:"name"`
	chefType  string `json:"chef_type"`
	jsonClass string `json:"json_class"`
	PublicKey string `json:"public_key"`
}

func (c *Chef) CreateClient(name string) (string, error) {
	payload := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(payload)
	err := encoder.Encode(map[string]string{"name": name})
	if err != nil {
		return "", err
	}

	responce, err := c.Post("clients", "application/json", nil, payload)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(responce.Body)
	responceHash := zhash.NewHash()
	json.Unmarshal(body, &responceHash)

	switch responce.StatusCode {
	case 201:
		key, err := responceHash.GetString("private_key")
		if err != nil {
			return "", errors.New(fmt.Sprintf("Error retrieving key: %s", err))
		}
		return key, nil
	default:
		errorMessage := getErrorMessage(responceHash)
		return "", errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s", responce.StatusCode, errorMessage))
	}
}

func (c *Chef) DeleteClient(name string) error {
	responce, err := c.Delete("clients/"+name, nil)
	if err != nil {
		return errors.New("chef-golang: " + err.Error())
	}

	body, err := ioutil.ReadAll(responce.Body)
	responceHash := zhash.NewHash()
	json.Unmarshal(body, &responceHash)

	switch responce.StatusCode {
	case 200:
		c.log.Debugf("Client %s deleted", name)
		return nil
	case 404:
		c.log.Noticef("Client %s not found", name)
		return nil
	default:
		errorMessage := getErrorMessage(responceHash)
		return errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s", responce.StatusCode, errorMessage))
	}

}
