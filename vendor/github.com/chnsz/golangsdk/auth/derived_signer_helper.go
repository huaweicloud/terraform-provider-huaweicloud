package auth

import (
	"net/http"

	"github.com/chnsz/golangsdk/auth/core/signer"
)

// SignDerived method is used to generate a derived authorization header and set it to the request header.
func SignDerived(request *http.Request, ak, sk, derivedAuthServiceName, regionID string) error {
	SignObj := signer.DerivedSigner{
		Key:                    ak,
		Secret:                 sk,
		DerivedAuthServiceName: derivedAuthServiceName,
		RegionID:               regionID,
	}
	return SignObj.Sign(request)
}
