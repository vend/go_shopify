package shopify

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

var api API
var remoteEnabled = false

func init() {
	if os.Getenv("SHOPIFY_API_PERM_TOKEN") != "" && os.Getenv("SHOPIFY_API_SHOP") != "" {
		remoteEnabled = true
		api = API{
			Shop:        os.Getenv("SHOPIFY_API_SHOP"),
			AccessToken: os.Getenv("SHOPIFY_API_PERM_TOKEN"),
		}
	} else if os.Getenv("SHOPIFY_API_TOKEN") != "" && os.Getenv("SHOPIFY_API_SECRET") != "" && os.Getenv("SHOPIFY_API_SHOP") != "" {
		remoteEnabled = true
		api = API{
			Shop:   os.Getenv("SHOPIFY_API_SHOP"),
			Token:  os.Getenv("SHOPIFY_API_TOKEN"),
			Secret: os.Getenv("SHOPIFY_API_SECRET"),
		}
	} else {
		log.Printf("Remote tests disabled, set SHOPIFY_API_KEY, SHOPIFY_API_SECRET, SHOPIFY_API_HOST, SHOPIFY_API_PERM_TOKEN")
	}
}

func TestReadProducts(t *testing.T) {
	if !remoteEnabled {
		return
	}

	products, err := api.Product(389374712)

	if err != nil {
		t.Errorf("Error fetching products: %v", err)
	}

	fmt.Printf("\n\nproducts are %#v\n\n", products)
}

func TestProductsCount(t *testing.T) {
	if !remoteEnabled {
		return
	}

	_, err := api.ProductsCount(nil)

	if err != nil {
		t.Errorf("Error fetching products count: %v", err)
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

func TestListCreateGetDeleteWebhook(t *testing.T) {
	if !remoteEnabled {
		return
	}

	// List and delete all
	_, err := api.Webhooks()

	if err != nil {
		fmt.Printf("Err fetching webhooks: %v", err)
	}

	// create
	newHook := api.NewWebhook()

	newHook.Address = "https://aaa.ngrok.com/service/hook"
	newHook.Format = "json"
	newHook.Topic = "orders/delete"
	err = newHook.Save(nil)
	if err != nil {
		t.Fatalf("Error creating webhook: %v", err)
	}

	//get
	hook, err := api.Webhook(newHook.Id)
	if err != nil {
		t.Errorf("Error fetching webhook (%v): %v", newHook.Id, err)
	}

	if hook.Id != newHook.Id {
		t.Errorf("Expected retrieved webhook to have the same ID as newly created webhook")
	}

	// clean up
	err = newHook.Delete()
	if err != nil {
		t.Errorf("Error deleting webhook: %s", err)
	}
}

func TestListCreateGetDeleteProduct(t *testing.T) {
	if !remoteEnabled {
		return
	}

	// List and delete all
	_, err := api.Products(&ProductsOptions{})
	if err != nil {
		fmt.Printf("Err fetching products: %v", err)
	}

	// create
	newProduct := api.NewProduct()
	newProduct.Title = "T-shirt"
	newProduct.PublishedAt = time.Now().String()
	newProduct.ProductType = "shirts"
	err = newProduct.Save(nil)
	if err != nil {
		t.Fatalf("Error saving product: %s", err)
	}
	if newProduct.ID == 0 {
		t.Errorf("Missing ID for newly created product")
	}

	// get new product by id
	product, err := api.Product(newProduct.ID)

	if err != nil {
		t.Errorf("Error fetching product (%v): %v", newProduct.ID, err)
	}

	if product.ID != newProduct.ID {
		t.Errorf("Expected retrieved product to have the same ID as newly created product")
	}

	// clean up
	err = newProduct.Delete()
	if err != nil {
		t.Errorf("Error deleting product: %s", err)
	}
}

func TestCreateWebhook(t *testing.T) {
	if !remoteEnabled {
		return
	}

	webhooks, err := api.Webhooks()

	if err != nil {
		fmt.Printf("Err fetching webhooks: %v", err)
	}

	for _, v := range webhooks {
		fmt.Printf("Existing webhook: %#v", v)
	}

	webhook := api.NewWebhook()

	webhook.Address = "https://aaa.ngrok.com/service/hook"
	webhook.Format = "json"
	webhook.Topic = "orders/delete"
	err = webhook.Save(nil)

	if err != nil {
		t.Errorf("Error creating webhook: %v", err)
	}

	fmt.Printf("\n\nwebhooks are %#v\n\n", webhook)
}

func TestNewProduct(t *testing.T) {
	if !remoteEnabled {
		return
	}

	product := api.NewProduct()
	product.Title = "T-shirt"
	product.PublishedAt = time.Now().String()
	product.ProductType = "shirts"
	err := product.Save(nil)
	if err != nil {
		t.Errorf("Error saving product: %s", err)
	}
	fmt.Printf("New product ID is: %d\n", product.ID)
}
