package shopify


import (
  
    "bytes"
  
    "encoding/json"
  
    "fmt"
  
    "time"
  
)


type Asset struct {
  
    Attachment string `json:"attachment"`

    ContentType string `json:"content_type"`
  
    CreatedAt time.Time `json:"created_at"`

    Key string `json:"key"`
  
    Size int64 `json:"size"`
  
    SourceKey string `json:"source_key"`
  
    Src string `json:"src"`

    ThemeId int64 `json:"theme_id"`
  
    UpdatedAt time.Time `json:"updated_at"`
  
    Value string `json:"value"`
  

  
    api *API
  
}


func (api *API) Assets(themeId int64) ([]Asset, error) {
  
  endpoint := fmt.Sprintf("/admin/themes/%d/assets.json", themeId)
  res, status, err := api.request(endpoint, "GET", nil, nil)

  if err != nil {
    return nil, err
  }

  if status != 200 {
    return nil, fmt.Errorf("Status returned: %d", status)
  }

  r := &map[string][]Asset{}
  err = json.NewDecoder(res).Decode(r)

  fmt.Printf("things are: %v\n\n", *r)

  result := (*r)["assets"]

	if err != nil {
		return nil, err
  }

  for _, v := range result {
    v.api = api
  }

  return result, nil
}




func (api *API) Asset(themeId int64, assetKey string) (*Asset, error) {
  endpoint := fmt.Sprintf("/admin/themes/%d/assets.json?asset=%s&theme_id=%d", themeId, assetKey, themeId )

  res, status, err := api.request(endpoint, "GET", nil, nil)

  if err != nil {
    return nil, err
  }

  if status != 200 {
    return nil, fmt.Errorf("Status returned: %d", status)
  }

  r := map[string]Asset{}
  err = json.NewDecoder(res).Decode(&r)

  fmt.Printf("things are: %v\n\n", r)

  result := r["asset"]

	if err != nil {
		return nil, err
  }

  result.api = api

  return &result, nil
}


func (api *API) NewAsset() *Asset {
  return &Asset{api: api}
}


func (obj *Asset) Save() (error) {
  endpoint := fmt.Sprintf("/admin/themes/%d/asset.json?asset=%s&theme_id=%d", obj.ThemeId, obj.Key, obj.ThemeId)
  method := "PUT"
  expectedStatus := 201

  body := map[string]*Asset{}
  body["asset"] = obj

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
    } else {
      return fmt.Errorf("Status %d, and error parsing body: %s", status, err)
    }
  }

  r := map[string]Asset{}
  err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
  }

  fmt.Printf("things are: %v\n\n", r)

  *obj = r["asset"]

  fmt.Printf("things are: %v\n\n", res)

  return nil
}





