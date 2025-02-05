// HWS API Gateway Signature
// based on https://github.com/datastream/aws/blob/master/signv4.go
// Copyright (c) 2014, Xianjie
// License that can be found in the LICENSE file

package signer

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/chnsz/golangsdk/auth/core/signer"
)

const (
	sdkHmacSha256     = "SDK-HMAC-SHA256"
	xSdkContentSha256 = "X-Sdk-Content-Sha256"
)

func Sign(req *http.Request, ak, sk string) (map[string]string, error) {
	return Signer{}.Sign(req, ak, sk)
}

type Signer struct {
}

// Sign SignRequest set Authorization header
func (s Signer) Sign(req *http.Request, ak, sk string) (map[string]string, error) {
	err := checkAKSK(ak, sk)
	if err != nil {
		return nil, err
	}

	processContentHeader(req, xSdkContentSha256)
	t := extractTime(req.Header.Get(signer.HeaderXDateTime))
	headerDate := t.UTC().Format(BasicDateFormat)
	req.Header.Set(signer.HeaderXDateTime, t.UTC().Format(BasicDateFormat))
	additionalHeaders := map[string]string{HeaderXDate: headerDate}

	signedHeaders := extractSignedHeaders(req)

	cr, err := canonicalRequest(req, signedHeaders, xSdkContentSha256, sha256HasherInst)
	if err != nil {
		return nil, err
	}

	sts, err := stringToSign(sdkHmacSha256, cr, t, sha256HasherInst)
	if err != nil {
		return nil, err
	}

	sig, err := s.signStringToSign(sts, []byte(sk))
	if err != nil {
		return nil, err
	}

	additionalHeaders[HeaderAuthorization] = authHeaderValue(sdkHmacSha256, sig, ak, signedHeaders)
	return additionalHeaders, nil
}

// signStringToSign Create the Signature.
func (Signer) signStringToSign(stringToSign string, signingKey []byte) (string, error) {
	hmac, err := sha256HasherInst.hmac([]byte(stringToSign), signingKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hmac), nil
}

const (
	BasicDateFormat     = "20060102T150405Z"
	HeaderXDate         = "X-Sdk-Date"
	HeaderHost          = "host"
	HeaderAuthorization = "Authorization"
)

func checkAKSK(ak, sk string) error {
	if ak == "" {
		return errors.New("ak is required in credentials")
	}
	if sk == "" {
		return errors.New("sk is required in credentials")
	}

	return nil
}

// stringToSign Create a "String to Sign".
func stringToSign(alg, canonicalRequest string, t time.Time, hasher iHasher) (string, error) {
	canonicalRequestHash, err := hasher.hashHexString([]byte(canonicalRequest))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%s", alg, t.UTC().Format(BasicDateFormat), canonicalRequestHash), nil
}

// authHeaderValue Get the finalized value for the "Authorization" header.
// The signature parameter is the output from stringToSign
func authHeaderValue(alg, sig, ak string, signedHeaders []string) string {
	return fmt.Sprintf("%s Access=%s, SignedHeaders=%s, Signature=%s",
		alg,
		ak,
		strings.Join(signedHeaders, ";"),
		sig)
}

func processContentHeader(req *http.Request, contentHeader string) {
	contentType := req.Header.Get("Content-Type")

	if contentType != "" && !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "application/bson") {
		req.Header.Set(contentHeader, "UNSIGNED-PAYLOAD")
	}
}

func canonicalRequest(req *http.Request, signedHeaders []string, contentHeader string, hasher iHasher) (string, error) {
	hexEncode, err := getContentHash(req, contentHeader, hasher)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		req.Method,
		canonicalURI(req),
		canonicalQueryString(req),
		canonicalHeaders(req, signedHeaders),
		strings.Join(signedHeaders, ";"), hexEncode), nil
}

func getContentHash(req *http.Request, contentHeader string, hasher iHasher) (string, error) {
	if content := req.Header.Get(contentHeader); content != "" {
		return content, nil
	}

	buffer, err := getBodyToBytes(req)
	if err != nil {
		return "", err
	}

	data := buffer.Bytes()
	hexEncode, err := hasher.hashHexString(data)
	if err != nil {
		return "", err
	}
	return hexEncode, nil
}

func extractTime(date string) time.Time {
	t, err := time.Parse(BasicDateFormat, date)
	if date == "" || err != nil {
		return time.Now()
	}
	return t
}

// canonicalURI returns request uri
func canonicalURI(r *http.Request) string {
	pattens := strings.Split(r.URL.Path, "/")

	var uri []string
	for _, v := range pattens {
		uri = append(uri, escape(v))
	}

	urlPath := strings.Join(uri, "/")
	if len(urlPath) == 0 || urlPath[len(urlPath)-1] != '/' {
		urlPath += "/"
	}

	return urlPath
}

func canonicalQueryString(r *http.Request) string {
	var keys []string
	query := r.URL.Query()
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var a []string
	for _, key := range keys {
		k := escape(key)
		sort.Strings(query[key])
		for _, v := range query[key] {
			kv := fmt.Sprintf("%s=%s", k, escape(v))
			a = append(a, kv)
		}
	}
	queryStr := strings.Join(a, "&")
	r.URL.RawQuery = queryStr

	return queryStr
}

func canonicalHeaders(r *http.Request, signerHeaders []string) string {
	var a []string
	header := make(map[string][]string)
	for k, v := range r.Header {
		header[strings.ToLower(k)] = v
	}

	for _, key := range signerHeaders {
		value := header[key]
		if strings.EqualFold(key, HeaderHost) {
			header[HeaderHost] = []string{r.Host}
		}

		sort.Strings(value)
		for _, v := range value {
			a = append(a, key+":"+strings.TrimSpace(v))
		}
	}

	return fmt.Sprintf("%s\n", strings.Join(a, "\n"))
}

func extractSignedHeaders(r *http.Request) []string {
	var sh []string
	for key := range r.Header {
		if strings.HasPrefix(strings.ToLower(key), "content-type") || strings.Contains(key, "_") {
			continue
		}
		sh = append(sh, strings.ToLower(key))
	}
	sort.Strings(sh)

	return sh
}

func shouldEscape(c byte) bool {
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '_' || c == '-' || c == '~' || c == '.' {
		return false
	}
	return true
}

func escape(s string) string {
	hexCount := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

func getBodyToBytes(req *http.Request) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return buf, err
	}

	if len(body) != 0 {
		v := reflect.ValueOf(req.Body)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Kind() == reflect.String {
			buf.WriteString(v.Interface().(string))
		} else {
			var err error
			contentType := req.Header.Get("Content-Type")

			switch contentType {
			case "application/xml":
				encoder := xml.NewEncoder(buf)
				err = encoder.Encode(req.Body)
				if err != nil {
					return nil, err
				}
			case "application/bson":
				buffer, err := bson.Marshal(req.Body)
				if err != nil {
					return nil, err
				}
				buf.Write(buffer)
			default:
				encoder := json.NewEncoder(buf)
				encoder.SetEscapeHTML(false)
				err = encoder.Encode(req.Body)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	return buf, nil
}
