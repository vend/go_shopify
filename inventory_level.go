package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
)

//InventoryLevel struct for Shopify inventory_level
type InventoryLevel struct {
	InventoryItemID     int64  `json:"inventory_item_id"`
	LocationID          int64  `json:"location_id"`
	Available           int64  `json:"available,omitempty"`
	AvailableAdjustment int64  `json:"available_adjustment,omitempty"`
	UpdatedAt           string `json:"updated_at,omitempty"`

	api *API
}

//Connect connect an inventory item to a location
func (obj *InventoryLevel) Connect() error {
	endpoint := fmt.Sprintf("/admin/inventory_levels/connect.json")
	method := "POST"
	return request(endpoint, method, obj)
}

//Set set an inventory level for a variant w. location id
func (obj *InventoryLevel) Set() error {
	endpoint := fmt.Sprintf("/admin/inventory_levels/set.json")
	method := "POST"
	return request(endpoint, method, obj)
}

//Adjust adjust an inventory level for a inventory item w. location id
func (obj *InventoryLevel) Adjust() error {
	endpoint := fmt.Sprintf("/admin/inventory_levels/adjust.json")
	method := "POST"
	return request(endpoint, method, obj)
}

//Delete delete an inventory level for a inventory item w. location id
func (obj *InventoryLevel) Delete() error {
	endpoint := fmt.Sprintf("/admin/inventory_levels.json?inventory_item_id=%&location_id=%", obj.InventoryItemID, obj.LocationID)
	method := "DELETE"
	return request(endpoint, method, obj)
}

func request(endpoint, method string, obj *InventoryLevel) error {
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
