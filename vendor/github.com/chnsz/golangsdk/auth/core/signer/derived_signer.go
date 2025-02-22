package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/hkdf"
)

const (
	DerivedDateFormat = "20060102"
	AlgorithmV11      = "V11-HMAC-SHA256"
)

// Derived algorithm structure, all fields are required
type DerivedSigner struct {
	Key                    string
	Secret                 string
	DerivedAuthServiceName string
	RegionID               string
}

// Sign Use derivative algorithms to authenticate AK/SK
func (s DerivedSigner) Sign(request *http.Request) error {
	if s.DerivedAuthServiceName == "" {
		return errors.New("DerivedAuthServiceName is required in credentials when using derived auth")
	}
	if s.RegionID == "" {
		return errors.New("RegionID is required in credentials when using derived auth")
	}

	var (
		t    time.Time
		err  error
		date string
	)

	if date = request.Header.Get(HeaderXDateTime); date != "" {
		t, err = time.Parse(DateFormat, date)
	}
	if err != nil || date == "" {
		t = time.Now()
		request.Header.Set(HeaderXDateTime, t.UTC().Format(DateFormat))
	}

	signedHeaders := SignedHeaders(request)
	canonicalRequest, err := CanonicalRequest(request, signedHeaders)
	if err != nil {
		return err
	}

	info := t.UTC().Format(DerivedDateFormat) + "/" + s.RegionID + "/" + s.DerivedAuthServiceName
	stringToSignStr, err := s.stringToSign(canonicalRequest, info, t)
	if err != nil {
		return err
	}

	derivationKey, err := s.getDerivationKey(info)
	if err != nil {
		return err
	}

	sig, err := s.signStringToSign(stringToSignStr, []byte(derivationKey))
	if err != nil {
		return err
	}

	authValueStr := s.authHeaderValue(sig, s.Key, info, signedHeaders)
	request.Header.Set(HeaderXAuthorization, authValueStr)

	return nil
}

func (s DerivedSigner) stringToSign(canonicalRequest, info string, t time.Time) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(canonicalRequest))
	if err != nil {
		return "", err
	}

	hashSum := hash.Sum(nil)
	canonicalRequestHash := hex.EncodeToString(hashSum)

	return fmt.Sprintf("%s\n%s\n%s\n%s", AlgorithmV11, t.UTC().Format(DateFormat), info, canonicalRequestHash), nil
}

func (s DerivedSigner) getDerivationKey(info string) (string, error) {
	hash := sha256.New
	derivationKeyReader := hkdf.New(hash, []byte(s.Secret), []byte(s.Key), []byte(info))
	derivationKey := make([]byte, 32)
	_, err := io.ReadFull(derivationKeyReader, derivationKey)
	return hex.EncodeToString(derivationKey), err
}

func (s DerivedSigner) signStringToSign(stringToSign string, signingKey []byte) (string, error) {
	hash := hmac.New(sha256.New, signingKey)
	if _, err := hash.Write([]byte(stringToSign)); err != nil {
		return "", err
	}
	hm := hash.Sum(nil)
	return fmt.Sprintf("%x", hm), nil
}

func (s DerivedSigner) authHeaderValue(signature, accessKey, info string, signedHeaders []string) string {
	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		AlgorithmV11,
		accessKey,
		info,
		strings.Join(signedHeaders, ";"),
		signature)
}
