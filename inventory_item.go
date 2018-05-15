package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// InventoryItem a struct to represent Shpoify' inventory_item.
type InventoryItem struct {
	ID        int64     `json:"id"`
	Sku       string    `json:"sku"`
	Tracked   bool      `json:"tracked"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	api *API
}

// InventoryItem Get one inventoryItem from api by inventory_item_id.
func (api *API) InventoryItem(id int64) (*InventoryItem, error) {
	endpoint := fmt.Sprintf("/admin/inventory_items/%d.json", id)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, newErrorResponse(status, nil, res)
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

// InventoryItems Get a list of inventoryItems from api, max 100 items.
func (api *API) InventoryItems() ([]InventoryItem, error) {
	res, status, err := api.request("/admin/inventory_items.json", "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, newErrorResponse(status, nil, res)
	}

	r := &struct {
		InventoryItems []InventoryItem `json:"inventory_items"`
	}{}
	err = json.NewDecoder(res).Decode(r)

	if err != nil {
		return nil, err
	}

	for _, v := range r.InventoryItems {
		v.api = api
	}

	return r.InventoryItems, nil
}

//Update update an existing inventory item based on inventory_item_id
func (obj *InventoryItem) Update() error {
	endpoint := fmt.Sprintf("/admin/inventory_items/%d.json", obj.ID)
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

	r := &struct {
		InventoryItem InventoryItem `json:"inventory_item"`
	}{}
	err = json.NewDecoder(res).Decode(r)
	if err != nil {
		return err
	}

	api := obj.api
	*obj = r.InventoryItem
	obj.api = api

	return nil
}
