package shopify

import "time"

//Config Simple struct for configuring diff vars of controlling the rate of making requests
type Config struct {
	BucketLimit        int           // default 30
	MaxRetries         int           // default 3
	MinBackoffValue    int           // minium backoff value
	MaxBackoffValue    int           // maxium backoff value
	MinBackOffTimeUnit time.Duration // Second, Millisecond etc.
	MaxBackOffTimeUnit time.Duration // Second, Millisecond etc.
	RefillRate         float64
}

//DefaultConfig return a default valued config
func DefaultConfig() Config {
	return Config{
		BucketLimit:        30,
		MaxRetries:         3,
		MinBackoffValue:    1,
		MaxBackoffValue:    4,
		MinBackOffTimeUnit: time.Second,
		MaxBackOffTimeUnit: time.Second,
		RefillRate:         float64(0.5),
	}
}
