package shopify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"strconv"
)

type Product struct {
	BodyHtml       string      `json:"body_html,omitempty"`
	CreatedAt      string      `json:"created_at,omitempty"`
	Handle         string      `json:"handle,omitempty"`
	ID             int64       `json:"id,omitempty"`
	Images         interface{} `json:"images,omitempty"`
	Options        []Option    `json:"options,omitempty"`
	ProductType    string      `json:"product_type,omitempty"`
	PublishedAt    string      `json:"published_at,omitempty"`
	PublishedScope string      `json:"published_scope,omitempty"`
	Tags           string      `json:"tags,omitempty"`
	TemplateSuffix string      `json:"template_suffix,omitempty"`
	Title          string      `json:"title,omitempty"`
	UpdatedAt      string      `json:"updated_at,omitempty"`
	Variants       []Variant   `json:"variants,omitempty"`
	Vendor         string      `json:"vendor,omitempty"`

	api *API
}

type ProductsOptions struct {
	IDs             string `url:"ids,omitempty"`
	Limit           int    `url:"limit,omitempty"`
	Page            int    `url:"page,omitempty"`
	SinceID         int64  `url:"since_id,omitempty"`
	Title           string `url:"title,omitempty"`
	Vendor          string `url:"vendor,omitempty"`
	Handle          string `url:"handle,omitempty"`
	ProductType     string `url:"product_type,omitempty"`
	CollectionID    string `url:"collection_id,omitempty"`
	CreatedAtMin    string `url:"created_at_min,omitempty"`
	CreatedAtMax    string `url:"created_at_max,omitempty"`
	UpdatedAtMin    string `url:"updated_at_min,omitempty"`
	UpdatedAtMax    string `url:"updated_at_max,omitempty"`
	PublishedAtMin  string `url:"published_at_min,omitempty"`
	PublishedAtMax  string `url:"published_at_max,omitempty"`
	PublishedStatus string `url:"published_status,omitempty"`
	Fields          string `url:"fields,omitempty"`
}

func (api *API) Products(options *ProductsOptions) ([]*Product, error) {

	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("/admin/products.json?%v", qs)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]*Product{}
	err = json.NewDecoder(res).Decode(r)
	if err != nil {
		return nil, err
	}

	result := (*r)["products"]
	for _, v := range result {
		v.api = api
	}

	return result, nil
}

type ProductsCountOptions struct {
	Vendor          string `url:"vendor,omitempty"`
	ProductType     string `url:"product_type,omitempty"`
	CollectionID    string `url:"collection_id,omitempty"`
	CreatedAtMin    string `url:"created_at_min,omitempty"`
	CreatedAtMax    string `url:"created_at_max,omitempty"`
	UpdatedAtMin    string `url:"updated_at_min,omitempty"`
	UpdatedAtMax    string `url:"updated_at_max,omitempty"`
	PublishedAtMin  string `url:"published_at_min,omitempty"`
	PublishedAtMax  string `url:"published_at_max,omitempty"`
	PublishedStatus string `url:"published_status,omitempty"`
}

func (api *API) ProductsCount(options *ProductsCountOptions) (int, error) {

	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("/admin/products/count.json?%v", qs)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return 0, err
	}

	if status != 200 {
		return 0, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]interface{}{}
	err = json.NewDecoder(res).Decode(&r)

	result, _ := strconv.Atoi(fmt.Sprintf("%v", r["count"]))
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (api *API) Product(id int64) (*Product, error) {
	endpoint := fmt.Sprintf("/admin/products/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Product{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["product"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewProduct() *Product {
	return &Product{api: api}
}

type ProductsMetafieldsOptions struct {
	Limit        int    `url:"limit,omitempty"`
	SinceID      string `url:"since_id,omitempty"`
	CreatedAtMin string `url:"created_at_min,omitempty"`
	CreatedAtMax string `url:"created_at_max,omitempty"`
	UpdatedAtMin string `url:"updated_at_min,omitempty"`
	UpdatedAtMax string `url:"updated_at_max,omitempty"`
	Namepace     string `url:"namepace,omitempty"`
	Key          string `url:"key,omitempty"`
	ValueType    string `url:"value_type,omitempty"`
	Fields       string `url:"fields,omitempty"`
}

func (obj *Product) Metafields(options *ProductsMetafieldsOptions) ([]*Metafield, error) {
	if obj == nil || obj.api == nil {
		return nil, errors.New("Product is nil")
	}
	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("/admin/products/%d/metafields.json?%v", obj.ID, qs)
	res, status, err := obj.api.request(endpoint, "GET", nil, nil)

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
		v.api = obj.api
	}

	return result, nil
}

//func (obj *Product) Save() error {
//	endpoint := fmt.Sprintf("/admin/products/%d.json", obj.Id)
//	method := "PUT"
//	expectedStatus := 200
//
//	if obj.Id == 0 {
//		endpoint = fmt.Sprintf("/admin/products.json")
//		method = "POST"
//		expectedStatus = 201
//	}
//
//	body := map[string]*Product{}
//	body["product"] = obj
//
//	buf := &bytes.Buffer{}
//	err := json.NewEncoder(buf).Encode(body)
//
//	if err != nil {
//		return err
//	}
//
//	res, status, err := obj.api.request(endpoint, method, nil, buf)
//
//	if err != nil {
//		return err
//	}
//
//	if status != expectedStatus {
//		r := errorResponse{}
//		err = json.NewDecoder(res).Decode(&r)
//		if err == nil {
//			return fmt.Errorf("Status %d: %v", status, r.Errors)
//		} else {
//			return fmt.Errorf("Status %d, and error parsing body: %s", status, err)
//		}
//	}
//
//	r := map[string]Product{}
//	err = json.NewDecoder(res).Decode(&r)
//
//	if err != nil {
//		return err
//	}
//
//	api := obj.api
//	*obj = r["product"]
//	obj.api = api
//
//	return nil
//}

func (obj *Product) Save(partial *Product) error {
	endpoint := fmt.Sprintf("/admin/products/%d.json", obj.ID)
	method := "PUT"
	expectedStatus := 200

	if obj.ID == 0 {
		endpoint = fmt.Sprintf("/admin/products.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*Product{}
	if partial == nil {
		body["product"] = obj
	} else {
		body["product"] = partial
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

	r := map[string]Product{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	api := obj.api
	*obj = r["product"]
	obj.api = api

	return nil
}

func (obj *Product) Delete() error {
	endpoint := fmt.Sprintf("/admin/products/%d.json", obj.ID)
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

func encodeOptions(v interface{}) string {
	str := ""
	qs, _ := query.Values(v)
	if qs != nil {
		str = qs.Encode()
	}
	return str
}
