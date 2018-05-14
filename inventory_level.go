package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type InventoryLevel struct {
	InventoryItemID int64  `json:"inventory_item_id"`
	LocationID      int64  `json:"location_id"`
	Available       int64  `json:"available,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`

	api *API
}

//Connect connect an inventory item to a location
func (obj *InventoryLevel) Connect() error {
	endpoint := fmt.Sprintf("/admin/inventory_levels/connect.json")
	method := "POST"
	expectedStatus := 200

	var buf bytes.Buffer
	body := map[string]*InventoryLevel{
		"inventory_level": obj,
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

	r := map[string]InventoryLevel{}
	err = json.NewDecoder(res).Decode(&r)
	if err != nil {
		return err
	}

	api := obj.api
	*obj = r["inventory_level"]
	obj.api = api
	return nil
}

//Set set an inventory level for a variant w. location id
func (obj *InventoryLevel) Set() error {
	endpoint := fmt.Sprintf("/admin/inventory_levels/set.json")
	method := "POST"
	expectedStatus := 200

	var buf bytes.Buffer
	body := map[string]*InventoryLevel{
		"inventory_level": obj,
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

	r := map[string]InventoryLevel{}
	err = json.NewDecoder(res).Decode(&r)
	if err != nil {
		return err
	}

	api := obj.api
	*obj = r["inventory_level"]
	obj.api = api

	return nil
}
