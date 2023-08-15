package auth

import (
	"net/http"

	"github.com/chnsz/golangsdk/auth/core/signer"
)

// Sign method is used to generate an authorization header and set it to the request header.
func Sign(request *http.Request, ak, sk string) error {
	SignObj := signer.Signer{
		Key:    ak,
		Secret: sk,
	}
	return SignObj.Sign(request)
}
