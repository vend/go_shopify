package shopify

type Option struct {
	ID        int64    `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Position  int64    `json:"position,omitempty"`
	ProductID int64    `json:"product_id,omitempty"`
	Values    []string `json:"values,omitempty"`
}
