package shopify

import "time"

//Config Simple struct for configuring diff vars of controlling the rate of making requests
type Config struct {
	BucketLimit     int           // default 30
	MaxRetries      int           // default 3
	MinBackoffValue time.Duration // minium backoff value
	MaxBackoffValue time.Duration // maxium backoff value
	RefillRate      float64
}

//DefaultConfig return a default valued config
func DefaultConfig() Config {
	return Config{
		BucketLimit:     30,
		MaxRetries:      3,
		MinBackoffValue: 1 * time.Second,
		MaxBackoffValue: 4 * time.Second,
		RefillRate:      float64(0.5),
	}
}
