package encryption

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/internal/vault/pgpkeys"
)

// RetrieveGPGKey returns the PGP key specified as the pgpKey parameter, or queries
// the public key from the keybase service if the parameter is a keybase username
// prefixed with the phrase "keybase:"
func RetrieveGPGKey(pgpKey string) (string, error) {
	const keybasePrefix = "keybase:"

	encryptionKey := pgpKey
	if strings.HasPrefix(pgpKey, keybasePrefix) {
		publicKeys, err := pgpkeys.FetchKeybasePubkeys([]string{pgpKey})
		if err != nil {
			return "", fmt.Errorf("Error retrieving Public Key for %s: %w", pgpKey, err)
		}
		encryptionKey = publicKeys[pgpKey]
	}

	return encryptionKey, nil
}

// EncryptValue encrypts the given value with the given encryption key. Description
// should be set such that errors return a meaningful user-facing response.
func EncryptValue(encryptionKey, value, description string) (fingerprint, encryptedValue string, err error) {
	fingerprints, encrypted, encryptErr := pgpkeys.EncryptShares([][]byte{[]byte(value)}, []string{encryptionKey})
	if encryptErr != nil {
		err = fmt.Errorf("Error encrypting %s: %w", description, encryptErr)
		return
	}
	if len(fingerprints) == 0 || len(encrypted) == 0 {
		err = errors.New("No PGP key fingerprints or encrypted values found")
		return
	}

	fingerprint = fingerprints[0]
	encryptedValue = base64.StdEncoding.EncodeToString(encrypted[0])
	return
}

// GetPGPFingerprint returns the fingerprint of the given base64-encoded PGP public key.
func GetPGPFingerprint(publicKey string) (string, error) {
	fingerprints, err := pgpkeys.GetFingerprints([]string{publicKey}, nil)
	if err != nil {
		return "", err
	}
	if len(fingerprints) == 0 {
		return "", errors.New("No PGP key fingerprints found")
	}

	return fingerprints[0], nil
}
