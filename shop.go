package shopify

import (
	"encoding/json"
	"fmt"
)

type Shop struct {
	Address1                        string  `json:"address1"`
	Address2                        string  `json:"address2"`
	City                            string  `json:"city"`
	Country                         string  `json:"country"`
	CountryCode                     string  `json:"country_code"`
	CountryName                     string  `json:"country_name"`
	CountyTaxes                     bool    `json:"county_taxes"`
	CreatedAt                       string  `json:"created_at"`
	Currency                        string  `json:"currency"`
	CustomerEmail                   string  `json:"customer_email"`
	Domain                          string  `json:"domain"`
	EligibleForCardReaderGiveaway   bool    `json:"eligible_for_card_reader_giveaway"`
	EligibleForPayments             bool    `json:"eligible_for_payments"`
	Email                           string  `json:"email"`
	ForceSsl                        bool    `json:"force_ssl"`
	GoogleAppsDomain                string  `json:"google_apps_domain"`
	GoogleAppsLoginEnabled          bool    `json:"google_apps_login_enabled"`
	HasDiscounts                    bool    `json:"has_discounts"`
	HasGiftCards                    bool    `json:"has_gift_cards"`
	HasStorefront                   bool    `json:"has_storefront"`
	IanaTimezone                    string  `json:"iana_timezone"`
	Id                              int64   `json:"id"`
	Latitude                        float64 `json:"latitude"`
	Longitude                       float64 `json:"longitude"`
	MoneyFormat                     string  `json:"money_format"`
	MoneyInEmailsFormat             string  `json:"money_in_emails_format"`
	MoneyWithCurrencyFormat         string  `json:"money_with_currency_format"`
	MoneyWithCurrencyInEmailsFormat string  `json:"money_with_currency_in_emails_format"`
	MultiLocationEnabled            bool    `json:"multi_location_enabled"`
	MyshopifyDomain                 string  `json:"myshopify_domain"`
	Name                            string  `json:"name"`
	PasswordEnabled                 bool    `json:"password_enabled"`
	Phone                           string  `json:"phone"`
	PlanDisplayName                 string  `json:"plan_display_name"`
	PlanName                        string  `json:"plan_name"`
	PrimaryLocale                   string  `json:"primary_locale"`
	PrimaryLocationId               int64   `json:"primary_location_id"`
	Province                        string  `json:"province"`
	ProvinceCode                    string  `json:"province_code"`
	RequiresExtraPaymentsAgreement  bool    `json:"requires_extra_payments_agreement"`
	SetupRequired                   bool    `json:"setup_required"`
	ShopOwner                       string  `json:"shop_owner"`
	Source                          string  `json:"source"`
	TaxShipping                     bool    `json:"tax_shipping"`
	TaxesIncluded                   bool    `json:"taxes_included"`
	Timezone                        string  `json:"timezone"`
	UpdatedAt                       string  `json:"updated_at"`
	Zip                             string  `json:"zip"`

	api *API
}

func (api *API) CurrentShop() (*Shop, error) {
	endpoint := "/admin/shop.json"

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Shop{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["shop"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}
