package chef

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"git.rn/devops/go-rpc.git"
	"git.rn/devops/zhash.git"
)

func NewEnvironment(name, description string, cookbookVersions map[string]string, attributes zhash.Hash) Environment {
	return Environment{
		Name:             name,
		Description:      description,
		Attributes:       attributes,
		CookbookVersions: cookbookVersions,
		JsonClass:        "Chef::Environment",
		ChefType:         "environment",
	}
}

type Environment struct {
	Name             string            `json:"name"`
	Attributes       zhash.Hash        `json:"attributes"`
	Description      string            `json:"description"`
	CookbookVersions map[string]string `json:"cookbook_versions"`
	JsonClass        string            `json:"json_class"`
	ChefType         string            `json:"chef_type"`
}

func (e Environment) String() string {
	pl, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "[Error marshaling environment: " + err.Error() + "]"
	}
	return string(pl)
}

func (c *Chef) CreateEnvironment(env Environment) error {
	payload, err := rpc.MarshalToJsonReader(env)
	if err != nil {
		return err
	}

	responce, err := c.Post("environments", nil, payload)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(responce.Body)
	responceHash := zhash.NewHash()
	json.Unmarshal(body, &responceHash)

	switch responce.StatusCode {
	case 201:
		return nil
	case 409:
		return errors.New(fmt.Sprintf("Environment %s already exists", env.Name))
	default:
		errorMessage := getErrorMessage(responceHash)
		return errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s; raw body: %s", responce.StatusCode, errorMessage, body))
	}
}

func (c *Chef) DeleteEnvironment(name string) error {
	responce, err := c.Delete("environments/"+name, nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(responce.Body)
	responceHash := zhash.NewHash()
	json.Unmarshal(body, &responceHash)

	switch responce.StatusCode {
	case 200:
		c.log.Debug("Environment %s deleted", name)
		return nil
	case 404:
		c.log.Notice("Environment %s not found", name)
		return nil
	default:
		errorMessage := getErrorMessage(responceHash)
		return errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s", responce.StatusCode, errorMessage))
	}
}