package shopify

import (
	"bytes"

	"encoding/json"

	"fmt"

	"time"
)

type CustomCollection struct {
	BodyHTML string `json:"body_html"`

	Handle string `json:"handle"`

	ID int64 `json:"id"`

	Image interface{} `json:"image,omitempty"`

	PublishedAt *time.Time `json:"published_at,omitempty"`

	PublishedScope string `json:"published_scope"`

	SortOrder string `json:"sort_order"`

	TemplateSuffix string `json:"template_suffix"`

	Title string `json:"title"`

	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	api *API
}

func (api *API) CustomCollections() ([]CustomCollection, error) {
	return api.CustomCollectionsWithOptions(&CollectionOptions{})
}

func (api *API) CustomCollectionsWithOptions(options *CollectionOptions) ([]CustomCollection, error) {
	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("/admin/custom_collections.json?%v", qs)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]CustomCollection{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["custom_collections"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) CustomCollection(id int64) (*CustomCollection, error) {
	endpoint := fmt.Sprintf("/admin/custom_collections/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]CustomCollection{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["custom_collection"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewCustomCollection() *CustomCollection {
	return &CustomCollection{api: api}
}

func (obj *CustomCollection) Save() error {
	endpoint := fmt.Sprintf("/admin/custom_collections/%d.json", obj.ID)
	method := "PUT"
	expectedStatus := 201

	if obj.ID == 0 {
		endpoint = fmt.Sprintf("/admin/custom_collections.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*CustomCollection{}
	body["custom_collection"] = obj

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
		r := errorResponse{}
		err = json.NewDecoder(res).Decode(&r)
		if err == nil {
			return fmt.Errorf("Status %d: %v", status, r.Errors)
		}

		return fmt.Errorf("Status %d, and error parsing body: %s", status, err)
	}

	r := map[string]CustomCollection{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	*obj = r["custom_collection"]

	return nil
}
