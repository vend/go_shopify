package shopify

import (
	"bytes"

	"encoding/json"

	"fmt"

	"time"
)

type Metafield struct {
	CreatedAt     time.Time `json:"created_at"`
	Description   string    `json:"description"`
	Id            int64     `json:"id"`
	Key           string    `json:"key"`
	Namespace     string    `json:"namespace"`
	OwnerId       int64     `json:"owner_id"`
	UpdatedAt     time.Time `json:"updated_at"`
	Value         string    `json:"value"`
	ValueType     string    `json:"value_type"`
	OwnerResource string    `json:"owner_resource"`
	api           *API
}

func (api *API) Metafields() ([]*Metafield, error) {
	res, status, err := api.request("/admin/metafields.json", "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string][]*Metafield{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["metafields"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) Metafield(id int64) (*Metafield, error) {
	endpoint := fmt.Sprintf("/admin/metafields/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Metafield{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["metafield"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewMetafield() *Metafield {
	return &Metafield{api: api}
}

func (obj *Metafield) Save() error {
	endpoint := fmt.Sprintf("/admin/metafields/%d.json", obj.Id)
	method := "PUT"
	expectedStatus := 201

	if obj.Id == 0 {
		endpoint = fmt.Sprintf("/admin/metafields.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*Metafield{}
	body["metafield"] = obj

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(body)

	if err != nil {
		return err
	}

	res, status, err := obj.api.request(endpoint, method, nil, buf)

	if err != nil {
		return err
	}

	if status != expectedStatus {
		return newErrorResponse(status, res)
	}

	r := map[string]Metafield{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	api := obj.api
	*obj = r["metafield"]
	obj.api = api

	return nil
}

func (obj *Metafield) SaveForProduct(productId int64) error {
	endpoint := fmt.Sprintf("/admin/products/%d/metafields/%d.json", productId, obj.Id)
	method := "PUT"
	expectedStatus := 200

	if obj.Id == 0 {
		endpoint = fmt.Sprintf("/admin/products/%d/metafields.json", productId)
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*Metafield{}
	body["metafield"] = obj

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(body)

	if err != nil {
		return err
	}

	res, status, err := obj.api.request(endpoint, method, nil, buf)

	if err != nil {
		return err
	}

	if status != expectedStatus {
		return newErrorResponse(status, res)
	}

	r := map[string]Metafield{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	api := obj.api
	*obj = r["metafield"]
	obj.api = api

	return nil
}
