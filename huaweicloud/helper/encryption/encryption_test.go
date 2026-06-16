package encryption

import (
	"os"
	"strings"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/internal/vault/pgpkeys"
)

type pgpTestConfig struct {
	publicKey   string
	privateKey  string
	fingerprint string
}

func loadPGPTestConfig(t *testing.T) pgpTestConfig {
	t.Helper()

	publicKey := strings.TrimSpace(os.Getenv("HW_PGP_PUBLIC_KEY"))
	if publicKey == "" {
		t.Skip("The PGP public key is not configured (HW_PGP_PUBLIC_KEY)")
	}

	privateKey := strings.TrimSpace(os.Getenv("HW_PGP_PRIVATE_KEY"))
	if privateKey == "" {
		t.Skip("The PGP private key is not configured (HW_PGP_PRIVATE_KEY)")
	}

	fingerprints, err := pgpkeys.GetFingerprints([]string{publicKey}, nil)
	if err != nil {
		t.Fatalf("Error computing PGP key fingerprint: %v", err)
	}
	if len(fingerprints) == 0 {
		t.Fatalf("No PGP key fingerprints found")
	}

	return pgpTestConfig{
		publicKey:   publicKey,
		privateKey:  privateKey,
		fingerprint: fingerprints[0],
	}
}

func TestRetrieveGPGKey_base64PublicKey(t *testing.T) {
	cfg := loadPGPTestConfig(t)

	key, err := RetrieveGPGKey(cfg.publicKey)
	if err != nil {
		t.Fatalf("Error retrieving PGP key: %v", err)
	}
	if key != cfg.publicKey {
		t.Fatalf("The retrieved PGP key is incorrect, want '%v', but got '%v'", cfg.publicKey, key)
	}
}

func TestRetrieveGPGKey_keybase(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Keybase API test in short mode")
	}

	key, err := RetrieveGPGKey("keybase:hashicorp")
	if err != nil {
		t.Fatalf("Error retrieving PGP key: %v", err)
	}
	if key == "" {
		t.Fatal("The retrieved PGP key is empty")
	}
}

func TestEncryptValue(t *testing.T) {
	cfg := loadPGPTestConfig(t)

	plaintext := "test-access-key-secret"
	gotFingerprint, encryptedValue, err := EncryptValue(cfg.publicKey, plaintext, "test secret")
	if err != nil {
		t.Fatalf("Error encrypting the secret: %v", err)
	}
	if gotFingerprint != cfg.fingerprint {
		t.Fatalf("The fingerprint of the PGP key used to encrypt the secret is incorrect, want '%v', but got '%v'", cfg.fingerprint, gotFingerprint)
	}
	if encryptedValue == "" {
		t.Fatal("The encrypted secret is empty")
	}

	decrypted, err := pgpkeys.DecryptBytes(encryptedValue, cfg.privateKey)
	if err != nil {
		t.Fatalf("Error decrypting the secret: %v", err)
	}
	if decrypted.String() != plaintext {
		t.Fatalf("The decrypted secret is incorrect, want '%v', but got '%v'", plaintext, decrypted.String())
	}
}
