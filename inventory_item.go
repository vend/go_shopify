package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
)

//InventoryItem a struct to reprensent Shpoify' inventory_item
type InventoryItem struct {
	ID        int64  `json:"id"`
	Sku       int64  `json:"sku"`
	Tracked   int64  `json:"tracked"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`

	api *API
}

//InventoryItem Get one inventoryItem from api by inventory_item_id
func (api *API) InventoryItem(id int64) (*InventoryItem, error) {
	endpoint := fmt.Sprintf("/admin/inventory_items/%.json", id)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string]InventoryItem{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["inventory_item"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

//InventoryItems Get a list of inventoryItems from api, max 100 items
func (api *API) InventoryItems() ([]InventoryItem, error) {
	res, status, err := api.request("/admin/inventory_items.json", "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]InventoryItem{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["inventory_items"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

//Update update an existing inventory item based on inventory_item_id
func (obj *InventoryItem) Update() error {
	endpoint := fmt.Sprintf("/admin/inventory_items/%.json", obj.ID)
	method := "PUT"
	expectedStatus := 200

	var buf bytes.Buffer
	body := map[string]*InventoryItem{
		"inventory_item": obj,
	}
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return err
	}
	reqBody := buf.Bytes()

	res, status, err := obj.api.request(endpoint, method, nil, &buf)
	if err != nil {
		return err
	}

	if status != expectedStatus {
		return newErrorResponse(status, reqBody, res)
	}

	r := map[string]InventoryItem{}
	err = json.NewDecoder(res).Decode(&r)
	if err != nil {
		return err
	}

	api := obj.api
	*obj = r["inventory_item"]
	obj.api = api

	return nil
}
