package shopify

import (
	"encoding/json"

	"fmt"

	"time"
)

type Collect struct {
	CollectionId int64 `json:"collection_id"`

	CreatedAt time.Time `json:"created_at"`

	Featured bool `json:"featured"`

	Id int64 `json:"id"`

	Position int64 `json:"position"`

	ProductId int64 `json:"product_id"`

	UpdatedAt time.Time `json:"updated_at"`

	SortValue string `json:"sort_value"`

	api *API
}

type CollectOptions struct {
	Limit int `url:"limit,omitempty"`
	Page  int `url:"page,omitempty"`
}

func (api *API) NewCollect() *Collect {
	return &Collect{api: api}
}

func (api *API) Collects() ([]Collect, error) {
	return api.CollectsWithOptions(&CollectOptions{})
}

func (api *API) CollectsWithOptions(options *CollectOptions) ([]Collect, error) {
	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("/admin/collects.json?%v", qs)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]Collect{}
	err = json.NewDecoder(res).Decode(r)
	result := (*r)["collects"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) Collect(id int64) (*Collect, error) {
	endpoint := fmt.Sprintf("/admin/collects/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Collect{}
	err = json.NewDecoder(res).Decode(&r)
	result := r["collect"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}
