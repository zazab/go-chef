package chef

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"git.rn/devops/go-rpc.git"
	"github.com/zazab/zhash"
)

type Node struct {
	Name            string     `json:"name"`
	ChefEnvironment string     `json:"chef_environment"`
	RunList         []string   `json:"run_list"`
	Normal          zhash.Hash `json:"normal"`
}

func (c *Chef) CreateNode(node Node) error {
	pl, err := rpc.MarshalToJsonReader(node)
	if err != nil {
		return err
	}

	responce, err := c.Post("nodes", nil, pl)
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
		c.log.Debug("Node %s deleted", name)
		return nil
	case 404:
		c.log.Notice("Node %s not found", name)
		return nil
	default:
		errorMessage := getErrorMessage(responceHash)
		return errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s", responce.StatusCode, errorMessage))
	}
}
