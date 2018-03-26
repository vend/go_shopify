package shopify

// Image represents a single image on a shopify product.
// Currently we push up all images as new on each push as we don't store
// a reference to the remote image entity's id.
type Image struct {
	Src string `json:"src,omitempty"`
}
