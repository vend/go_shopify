package shopify

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jpillora/backoff"
	"github.com/vend/log"
)

const REFILL_RATE = float64(0.5) // 2 per second
const MAX_RETRIES = 3

type API struct {
	Shop        string // for e.g. demo-3.myshopify.com
	AccessToken string // permanent store access token
	Token       string // API client token
	Secret      string // API client secret for this application

	callLimit  int
	callsMade  int
	backoff    *backoff.Backoff
	retryCount int
}

// ErrorResponse is returned when an unexpected HTTP status code is received.
type ErrorResponse struct {
	// Errors is either a map of errors or a single error string returned
	// in the HTTP response, if it was possible to decode the response as JSON.
	Errors interface{} `json:"errors"`
	// StatusCode is the HTTP status code served in the response.
	StatusCode int `json:"-"`
	// ReqBody is the HTTP request body.
	ReqBody []byte `json:"-"`
	// Body is the HTTP response body.
	Body []byte `json:"-"`
	// BodyErr is any encoding/json error that occurred while trying to
	// unmarshal into the Errors field.
	BodyErr error `json:"-"`
}

func newErrorResponse(status int, reqBody []byte, body *bytes.Buffer) error {
	var r ErrorResponse
	r.StatusCode = status
	r.ReqBody = reqBody
	r.Body = body.Bytes()
	r.BodyErr = json.NewDecoder(body).Decode(&r)
	return &r
}

func (e *ErrorResponse) Error() string {
	ret := fmt.Sprintf("status %d: %v", e.StatusCode, e.Errors)
	if len(e.ReqBody) > 0 {
		ret += "; request body: " + string(e.ReqBody)
	}
	if e.BodyErr != nil {
		ret += "; error parsing body: " + e.BodyErr.Error()
	}
	if len(e.Body) > 0 {
		ret += "; response body: " + string(e.Body)
	}
	return ret
}

// Temporary returns true when the status code indicates that an error is probably
// temporary.
func (e *ErrorResponse) Temporary() bool {
	return e.StatusCode >= 500 || e.StatusCode == http.StatusTooManyRequests
}

func (api *API) request(endpoint string, method string, params map[string]interface{}, body *bytes.Buffer) (result *bytes.Buffer, status int, err error) {
	bucketLimit, err := strconv.Atoi(os.Getenv("BUCKET_LIMIT"))
	if err != nil {
		bucketLimit = 30
	}

	if api.backoff == nil {
		minBackoffSecond, err := strconv.ParseInt(os.Getenv("MIN_BACKOFF_SECOND"), 10, 64)
		if err != nil {
			minBackoffSecond = 1
		}
		maxBackoffSecond, err := strconv.ParseInt(os.Getenv("MAX_BACKOFF_SECOND"), 10, 64)
		if err != nil {
			maxBackoffSecond = 4
		}
		api.backoff = &backoff.Backoff{
			Min:    time.Duration(minBackoffSecond) * time.Second,
			Max:    time.Duration(maxBackoffSecond) * time.Second,
			Jitter: true,
		}
	}
	if api.callLimit == 0 {
		api.callLimit = bucketLimit
	}

	// Keep a copy of body so that we can use it when retrying.
	var bodyBackup *bytes.Buffer
	if body != nil {
		bodyBackup = new(bytes.Buffer)
		*bodyBackup = *body
	}

	uri := fmt.Sprintf("https://%s%s", api.Shop, endpoint)
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return
	}

	if api.Secret == "" {
		req.Header.Set("X-Shopify-Access-Token", api.AccessToken)
	} else {
		sum := md5.Sum([]byte(api.Secret + api.AccessToken))
		hexSum := hex.EncodeToString(sum[:])
		req.SetBasicAuth(api.Token, hexSum)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	calls, total := parseAPICallLimit(resp.Header.Get("HTTP_X_SHOPIFY_SHOP_API_CALL_LIMIT"))
	api.callsMade = calls
	api.callLimit = total

	status = resp.StatusCode
	if status == 429 { // statusTooManyRequests
		if api.retryCount < MAX_RETRIES {
			api.retryCount = api.retryCount + 1
			b := api.backoff.Duration()
			log.Global().WithField("backoff duration values", b).WithField("calls made ", calls).WithField("calls limit ", total).Error("time to sleep")
			time.Sleep(b)
			// try again
			return api.request(endpoint, method, params, bodyBackup)
		}
		// else just return
	}

	result = &bytes.Buffer{}
	if _, err = io.Copy(result, resp.Body); err != nil {
		return
	}
	return
}

func parseAPICallLimit(str string) (int, int) {
	tokens := strings.Split(str, "/")
	if len(tokens) != 2 {
		return 0, 0
	}
	calls, err := strconv.Atoi(tokens[0])
	if err != nil {
		return 0, 0
	}
	total, err := strconv.Atoi(tokens[1])
	if err != nil {
		return 0, 0
	}
	return calls, total
}
