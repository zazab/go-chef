package chef

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"git.rn/devops/go-rpc.git"
	"git.rn/devops/zhash.git"
)

type DataBagItem struct {
	Id   string     `json:"id"`
	Data zhash.Hash `json:"data"`
}

func (c *Chef) GetDatabagItemList(databag string) ([]string, error) {
	responce, err := c.Get("data/" + databag)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(responce.Body)
	if err != nil {
		return nil, err
	}
	databags := map[string]string{}
	err = json.Unmarshal(body, &databags)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for name, _ := range databags {
		result = append(result, name)
	}

	return result, err
}

func (c *Chef) GetDatabagItem(databag, item string) (DataBagItem, error) {
	responce, err := c.Get(fmt.Sprintf("data/%s/%s", databag, item))
	if err != nil {
		return DataBagItem{}, err
	}

	body, err := ioutil.ReadAll(responce.Body)
	if err != nil {
		return DataBagItem{}, err
	}
	c.log.Debug("Databag item body: %s", body)
	data := DataBagItem{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return DataBagItem{}, err
	}

	c.log.Debug("Databag=%[1]#v\nString=%[1]s", data)
	switch responce.StatusCode {
	case 200:
		return data, nil
	case 401:
		return DataBagItem{}, errors.New(fmt.Sprintf("Unauthorized. The user "+
			"which made the request is not authorized to perform the action. "+
			"Response: %s", responce))
	case 403:
		return DataBagItem{}, errors.New(fmt.Sprintf("Forbidden. The user which "+
			"made the request is not authorized to perform the action. "+
			"Response: %s", responce))
	case 404:
		return DataBagItem{}, errors.New(fmt.Sprintf("Requested databag item "+
			"%s/%s not found", databag, item))
	default:
		return DataBagItem{}, errors.New(fmt.Sprintf("Unknown response status code %d. "+
			"Response: %s", responce.StatusCode, responce))
	}
}

func (c *Chef) CreateDatabag(databag string) error {
	data := map[string]string{"name": databag}
	payload, err := rpc.MarshalToJsonReader(data)
	if err != nil {
		return err
	}

	responce, err := c.Post("data", nil, payload)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(responce.Body)
	responceHash := zhash.NewHash()
	json.Unmarshal(body, &responceHash)
	switch responce.StatusCode {
	case 201:
		return err
	case 409:
		return errors.New(fmt.Sprintf("Databag %s already exists", databag))
	default:
		errorMessage := getErrorMessage(responceHash)
		return errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s", responce.StatusCode, errorMessage))
	}
}

func (c *Chef) CreateDatabagItem(databag, item string, value zhash.Hash) error {
	payload, err := rpc.MarshalToJsonReader(
		DataBagItem{Id: item, Data: value})
	if err != nil {
		return err
	}

	responce, err := c.Post("data/"+databag, nil, payload)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(responce.Body)
	responceHash := zhash.NewHash()
	err = json.Unmarshal(body, &responceHash)
	switch responce.StatusCode {
	case 200:
		return err
	case 409:
		return errors.New(fmt.Sprintf("Databag item %s/%s already exists",
			databag, item))
	default:
		errorMessage := getErrorMessage(responceHash)
		return errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s", responce.StatusCode, errorMessage))
	}

	return nil
}

func (c *Chef) DeleteDatabag(databag string) error {
	responce, err := c.Delete("data/"+databag, nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(responce.Body)
	responceHash := zhash.NewHash()
	json.Unmarshal(body, &responceHash)

	switch responce.StatusCode {
	case 200:
		c.log.Debug("Databag %s deleted", databag)
		return nil
	case 404:
		c.log.Notice("Databag %s not found", databag)
		return nil
	default:
		errorMessage := getErrorMessage(responceHash)
		return errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s", responce.StatusCode, errorMessage))
	}
}

func (c *Chef) DeleteDatabagItem(databag, item string) error {
	responce, err := c.Delete(fmt.Sprintf("data/%s/%s", databag, item), nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(responce.Body)
	responceHash := zhash.NewHash()
	json.Unmarshal(body, &responceHash)

	switch responce.StatusCode {
	case 200:
		c.log.Debug("Databag item %s/%s deleted", databag, item)
		return nil
	case 404:
		c.log.Notice("Databag item %s/%s not found", databag, item)
		return nil
	default:
		errorMessage := getErrorMessage(responceHash)
		return errors.New(fmt.Sprintf("Response status code %d. "+
			"Error: %s", responce.StatusCode, errorMessage))
	}
}
