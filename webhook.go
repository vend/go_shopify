package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Webhook struct {
	Address             string        `json:"address,omitempty"`
	CreatedAt           time.Time     `json:"created_at,omitempty"`
	Fields              []interface{} `json:"fields,omitempty"`
	Format              string        `json:"format,omitempty"`
	Id                  int64         `json:"id,omitempty"`
	MetafieldNamespaces []interface{} `json:"metafield_namespaces,omitempty"`
	Topic               string        `json:"topic,omitempty"`
	UpdatedAt           time.Time     `json:"updated_at,omitempty"`
	api                 *API
}

func (api *API) Webhooks() ([]*Webhook, error) {
	res, status, err := api.request("/admin/webhooks.json", "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]*Webhook{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["webhooks"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) Webhook(id int64) (*Webhook, error) {
	endpoint := fmt.Sprintf("/admin/webhooks/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Webhook{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["webhook"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewWebhook() *Webhook {
	return &Webhook{api: api}
}

func (obj *Webhook) Save(partial *Webhook) error {
	endpoint := fmt.Sprintf("/admin/webhooks/%d.json", obj.Id)
	method := "PUT"
	expectedStatus := 200

	if obj.Id == 0 {
		endpoint = fmt.Sprintf("/admin/webhooks.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*Webhook{}
	if partial == nil {
		body["webhook"] = obj
	} else {
		body["webhook"] = partial
	}

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
		return newErrorResponse(status, nil, res)
	}

	r := map[string]Webhook{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	api := obj.api
	*obj = r["webhook"]
	obj.api = api

	return nil
}

func (obj *Webhook) Delete() error {
	endpoint := fmt.Sprintf("/admin/webhooks/%d.json", obj.Id)
	method := "DELETE"
	expectedStatus := 200

	res, status, err := obj.api.request(endpoint, method, nil, nil)

	if err != nil {
		return err
	}

	if status != expectedStatus {
		return newErrorResponse(status, nil, res)
	}

	return nil
}
