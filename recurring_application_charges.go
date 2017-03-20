package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	StatusAccepted = "accepted"
	StatusDeclined = "declined"
)

type RecurringApplicationCharge struct {
	ActivatedOn        string `json:"activated_on,omitempty"`
	APIClientID        int64  `json:"api_client_id,omitempty"`
	BillingOn          string `json:"billing_on,omitempty"`
	CancelledOn        string `json:"cancelled_on,omitempty"`
	ConfirmationURL    string `json:"confirmation_url,omitempty"`
	DecoratedReturnURL string `json:"decorated_return_url,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	ID                 int64  `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	Status             string `json:"status,omitempty"`
	Price              string `json:"price,omitempty"`
	ReturnURL          string `json:"return_url,omitempty"`
	Test               bool   `json:"test,omitempty"`
	TrialDays          int    `json:"trial_days,omitempty"`
	TrialEndsOn        string `json:"trial_ends_on,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`

	api *API
}

type RecurringApplicationChargeOptions struct {
	SinceID int64  `url:"since_id,omitempty"`
	Fields  string `url:"fields,omitempty"`
}

// RecurringApplicationCharges Retrieve all recurring application charges
func (api *API) RecurringApplicationCharges(options *RecurringApplicationChargeOptions) ([]*RecurringApplicationCharge, error) {

	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("/admin/recurring_application_charges.json?%v", qs)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]*RecurringApplicationCharge{}
	err = json.NewDecoder(res).Decode(r)
	if err != nil {
		return nil, err
	}

	result := (*r)["recurring_application_charges"]
	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) RecurringApplicationCharge(id int64) (*RecurringApplicationCharge, error) {
	endpoint := fmt.Sprintf("/admin/recurring_application_charges/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]RecurringApplicationCharge{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["recurring_application_charge"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewRecurringApplicationCharge() *RecurringApplicationCharge {
	return &RecurringApplicationCharge{api: api}
}

func (obj *RecurringApplicationCharge) Save() error {

	endpoint := fmt.Sprintf("/admin/recurring_application_charges.json")
	method := "POST"
	expectedStatus := 201

	body := map[string]*RecurringApplicationCharge{}
	body["recurring_application_charge"] = obj

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

	r := map[string]RecurringApplicationCharge{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	api := obj.api
	*obj = r["recurring_application_charge"]
	obj.api = api

	return nil
}

func (obj *RecurringApplicationCharge) Activate() error {
	endpoint := fmt.Sprintf("/admin/recurring_application_charges/%d/activate.json", obj.ID)
	method := "POST"
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

func (obj *RecurringApplicationCharge) Delete() error {
	endpoint := fmt.Sprintf("/admin/recurring_application_charges/%d.json", obj.ID)
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
