package chef

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/zazab/zhash"
)

type Node struct {
	Name            string     `json:"name"`
	ChefEnvironment string     `json:"chef_environment"`
	RunList         []string   `json:"run_list"`
	Normal          zhash.Hash `json:"normal"`
}

func (c *Chef) CreateNode(node Node) error {
	payload := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(payload)
	err := encoder.Encode(node)
	if err != nil {
		return err
	}

	responce, err := c.Post("nodes", "application/json", nil, payload)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(responce.Body)
	responceHash := zhash.NewHash()
	json.Unmarshal(body, &responceHash)

	switch responce.StatusCode {
	case 201:
		return nil
	default:
		errorMessage := getErrorMessage(responceHash)
		return errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s", responce.StatusCode, errorMessage))
	}
}

func (c *Chef) DeleteNode(name string) error {
	responce, err := c.Delete("nodes/"+name, nil)
	if err != nil {
		return errors.New("chef-golang: " + err.Error())
	}

	body, err := ioutil.ReadAll(responce.Body)
	responceHash := zhash.NewHash()
	json.Unmarshal(body, &responceHash)

	switch responce.StatusCode {
	case 200:
		c.log.Debugf("Node %s deleted", name)
		return nil
	case 404:
		c.log.Noticef("Node %s not found", name)
		return nil
	default:
		errorMessage := getErrorMessage(responceHash)
		return errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s", responce.StatusCode, errorMessage))
	}
}
