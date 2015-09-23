package chef

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"git.rn/devops/go-rpc.git"
	"github.com/zazab/zhash"
)

type Client struct {
	Name      string `json:"name"`
	chefType  string `json:"chef_type"`
	jsonClass string `json:"json_class"`
	PublicKey string `json:"public_key"`
}

func (c *Chef) CreateClient(name string) (string, error) {
	pl, err := rpc.MarshalToJsonReader(map[string]string{"name": name})
	if err != nil {
		return "", err
	}

	responce, err := c.Post("clients", nil, pl)
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
		return err
	}

	body, err := ioutil.ReadAll(responce.Body)
	responceHash := zhash.NewHash()
	json.Unmarshal(body, &responceHash)

	switch responce.StatusCode {
	case 200:
		c.log.Debug("Client %s deleted", name)
		return nil
	case 404:
		c.log.Notice("Client %s not found", name)
		return nil
	default:
		errorMessage := getErrorMessage(responceHash)
		return errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s", responce.StatusCode, errorMessage))
	}

}
