package shopee

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Signature struct {
	partnerID  int64
	partnerKey string
	apiVersion string
}

func NewSignature(partnerID int64, partnerKey, apiVersion string) *Signature {
	if apiVersion == "" {
		apiVersion = "v2"
	}
	return &Signature{partnerID: partnerID, partnerKey: partnerKey, apiVersion: apiVersion}
}

func (s *Signature) APIPath(path string) string {
	if strings.HasPrefix(path, "/api/") {
		return path
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "/api/" + s.apiVersion + path
}

func (s *Signature) Authorization(path string, timestamp int64) string {
	apiPath := s.APIPath(path)
	return s.sign(fmt.Sprintf("%d%s%d", s.partnerID, apiPath, timestamp))
}

func (s *Signature) Request(path string, timestamp int64, accessToken, shopID, merchantID string) string {
	apiPath := s.APIPath(path)
	base := fmt.Sprintf("%d%s%d", s.partnerID, apiPath, timestamp)
	if accessToken != "" {
		base += accessToken
	}
	if shopID != "" {
		base += shopID
	} else if merchantID != "" {
		base += merchantID
	}
	return s.sign(base)
}

func (s *Signature) sign(payload string) string {
	mac := hmac.New(sha256.New, []byte(s.partnerKey))
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
