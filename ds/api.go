package ds

import (
	"encoding/json"

	"github.com/silentrc/toolbox"
)

func (c *Cli) ListModel() ([]Data, error) {
	resp, err := toolbox.NewUtils().
		NewHttpUtils().
		Client(true).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.key).
		Get(c.url + "/models")
	if err != nil {
		return nil, err
	}
	list := ListModelResponse{}
	err = json.Unmarshal(resp.Body(), &list)
	if err != nil {
		return nil, err
	}
	return list.Data, nil
}

type ListModelResponse struct {
	Object string `json:"object"`
	Data   []Data `json:"data"`
}

type Data struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}
