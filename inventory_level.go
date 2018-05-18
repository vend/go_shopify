package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// InventoryLevel struct for Shopify inventory_level.
type InventoryLevel struct {
	InventoryItemID     int64     `json:"inventory_item_id"`
	LocationID          int64     `json:"location_id"`
	Available           int64     `json:"available,omitempty"`
	AvailableAdjustment int64     `json:"available_adjustment,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`

	api *API
}

// Connect connects an inventory item to a location.
func (obj *InventoryLevel) Connect() error {
	return requestInvLevel("/admin/inventory_levels/connect.json", "POST", obj)
}

// Set sets an inventory level for a variant w. location id.
func (obj *InventoryLevel) Set() error {
	return requestInvLevel("/admin/inventory_levels/set.json", "POST", obj)
}

// Adjust adjust an inventory level for a inventory item w. location id.
func (obj *InventoryLevel) Adjust() error {
	return requestInvLevel("/admin/inventory_levels/adjust.json", "POST", obj)
}

// Delete delete an inventory level for a inventory item w. location id.
func (obj *InventoryLevel) Delete() error {
	endpoint := fmt.Sprintf("/admin/inventory_levels.json?inventory_item_id=%d&location_id=%d", obj.InventoryItemID, obj.LocationID)
	expectedStatus := 204
	res, status, err := obj.api.request(endpoint, "DELETE", nil, nil)
	if err != nil {
		return err
	}
	if status != expectedStatus {
		return newErrorResponse(status, nil, res)
	}

	return nil
}

// requestInvLevel private func to make requests for inventory level.
func requestInvLevel(endpoint, method string, obj *InventoryLevel) error {
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

	r := &struct {
		InventoryLevel InventoryLevel `json:"inventory_level"`
	}{}
	err = json.NewDecoder(res).Decode(r)
	if err != nil {
		return err
	}

	api := obj.api
	*obj = r.InventoryLevel
	obj.api = api

	return nil
}
