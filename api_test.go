package shopify

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var api API
var remoteEnabled = false

func init() {

	if os.Getenv("SHOPIFY_API_TOKEN") != "" && os.Getenv("SHOPIFY_API_SECRET") != "" && os.Getenv("SHOPIFY_API_HOST") != "" {
		remoteEnabled = true
		api = API{
			URI:    os.Getenv("SHOPIFY_API_HOST"),
			Token:  os.Getenv("SHOPIFY_API_TOKEN"),
			Secret: os.Getenv("SHOPIFY_API_SECRET"),
		}
	} else {
		log.Printf("Remote tests disabled, set SHOPIFY_API_KEY, SHOPIFY_API_SECRET, SHOPIFY_API_HOST")
	}
}

func TestGetAssets(t *testing.T) {
	if !remoteEnabled {
		return
	}

	assets, err := api.Assets(122410883)

	if err != nil {
		t.Errorf("Error fetching assets: %v", err)
	}

	fmt.Printf("\n\assets are %#v\n\n", assets)
}

