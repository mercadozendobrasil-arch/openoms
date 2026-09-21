package shopee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	cfg Config
	sig *Signature
}

type Request struct {
	Method      string
	Path        string
	Query       url.Values
	Body        any
	AccessToken string
	ShopID      string
	MerchantID  string
	Timestamp   int64
}

type Envelope[T any] struct {
	Error     string `json:"error,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Warning   string `json:"warning,omitempty"`
	Response  *T     `json:"response,omitempty"`
}

type APIError struct {
	Code      string
	Message   string
	RequestID string
	Status    int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("shopee: %s: %s", e.Code, e.Message)
}

func NewClient(cfg Config) *Client {
	cfg = cfg.withDefaults()
	return &Client{cfg: cfg, sig: NewSignature(cfg.PartnerID, cfg.PartnerKey, cfg.APIVersion)}
}

func (c *Client) AuthorizationURL(redirectURL, state string, timestamp int64) (string, error) {
	if timestamp == 0 {
		timestamp = c.cfg.Now().Unix()
	}
	path := c.sig.APIPath("/shop/auth_partner")
	u, err := url.Parse(strings.TrimRight(c.cfg.BaseURL, "/") + path)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("partner_id", strconv.FormatInt(c.cfg.PartnerID, 10))
	q.Set("timestamp", strconv.FormatInt(timestamp, 10))
	q.Set("redirect", redirectURL)
	q.Set("sign", c.sig.Authorization(path, timestamp))
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (c *Client) Do(ctx context.Context, req Request, out any) error {
	if req.Timestamp == 0 {
		req.Timestamp = c.cfg.Now().Unix()
	}
	path := c.sig.APIPath(req.Path)
	u, err := url.Parse(strings.TrimRight(c.cfg.BaseURL, "/") + path)
	if err != nil {
		return err
	}
	q := req.Query
	if q == nil {
		q = make(url.Values)
	}
	q.Set("partner_id", strconv.FormatInt(c.cfg.PartnerID, 10))
	q.Set("timestamp", strconv.FormatInt(req.Timestamp, 10))
	if req.AccessToken != "" {
		q.Set("access_token", req.AccessToken)
	}
	if req.ShopID != "" {
		q.Set("shop_id", req.ShopID)
	}
	if req.MerchantID != "" {
		q.Set("merchant_id", req.MerchantID)
	}
	q.Set("sign", c.sig.Request(path, req.Timestamp, req.AccessToken, req.ShopID, req.MerchantID))
	u.RawQuery = q.Encode()

	var body io.Reader
	if req.Body != nil {
		raw, err := json.Marshal(req.Body)
		if err != nil {
			return err
		}
		body = bytes.NewReader(raw)
	}
	method := req.Method
	if method == "" {
		method = http.MethodPost
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Accept", "application/json")
	if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.cfg.HTTPClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var env struct {
		Error     string          `json:"error"`
		Message   string          `json:"message"`
		RequestID string          `json:"request_id"`
		Response  json.RawMessage `json:"response"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return fmt.Errorf("shopee: invalid JSON response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || env.Error != "" {
		code := env.Error
		if code == "" {
			code = resp.Status
		}
		return &APIError{Code: code, Message: env.Message, RequestID: env.RequestID, Status: resp.StatusCode}
	}
	if out == nil {
		return nil
	}
	if len(env.Response) > 0 && string(env.Response) != "null" {
		return json.Unmarshal(env.Response, out)
	}
	// Auth endpoints return token fields at the top level.
	return json.Unmarshal(raw, out)
}
