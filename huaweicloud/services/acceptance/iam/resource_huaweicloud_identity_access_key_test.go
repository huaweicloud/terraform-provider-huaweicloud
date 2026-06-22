package iam

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keybase/go-crypto/openpgp"
	"github.com/keybase/go-crypto/openpgp/packet"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/credentials"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getV3AccessKeyResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	iamClient, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return credentials.Get(iamClient, state.Primary.ID).Extract()
}

func TestAccV3AccessKey_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_access_key.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getV3AccessKeyResourceFunc)

		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3AccessKey_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|(\.\d{6}Z))$`)),
					resource.TestCheckNoResourceAttr(resourceName, "secret"),
				),
			},
			{
				Config: testAccV3AccessKey_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "inactive"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|(\.\d{6}Z))$`)),
				),
			},
		},
	})
}

func testAccV3AccessKey_basic_base(name string) string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}

resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = random_string.test.result
  enabled     = true
  description = "Created by terraform script"
}
`, name)
}

func testAccV3AccessKey_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_access_key" "test" {
  user_id     = huaweicloud_identity_user.test.id
  description = "Created by terraform script"
  secret_file = abspath("./credentials.csv")

  # Check the mode of the credential file (created by huaweicloud_identity_access_key resource and with a default name)
  # after the resource is created, want 600 access mode.
  provisioner "local-exec" {
    when    = create
    command = <<-EOT
      f="${abspath("./credentials.csv")}"
      perms=$(stat -c '%%a' "$f")
      if [ "$perms" != "600" ]; then
        echo "ERROR: $f has mode $perms, expected 600" >&2
        exit 1
      fi
    EOT
  }

  # Clean up the credential file (created by huaweicloud_identity_access_key resource and with a default name) after the
  # test is completed.
  provisioner "local-exec" {
    command = format("rm -f %%s", abspath("./credentials.csv"))
    when    = destroy
  }
}
`, testAccV3AccessKey_basic_base(name))
}

func testAccV3AccessKey_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_access_key" "test" {
  user_id     = huaweicloud_identity_user.test.id
  secret_file = abspath("./credentials.csv")
  status      = "inactive"

  # Check the mode of the credential file (created by huaweicloud_identity_access_key resource and with a default name)
  # after the resource is created, want 600 access mode.
  provisioner "local-exec" {
    when    = create
    command = <<-EOT
      f="${abspath("./credentials.csv")}"
      perms=$(stat -c '%%a' "$f")
      if [ "$perms" != "600" ]; then
        echo "ERROR: $f has mode $perms, expected 600" >&2
        exit 1
      fi
    EOT
  }

  # Clean up the credential file (created by huaweicloud_identity_access_key resource and with a default name) after the
  # test is completed.
  provisioner "local-exec" {
    command = format("rm -f %%s", abspath("./credentials.csv"))
    when    = destroy
  }
}
`, testAccV3AccessKey_basic_base(name))
}

// defaultV3AccessKeyCredentialFilePath returns the default credentials CSV path or the path specified by secret_file.
func defaultV3AccessKeyCredentialFilePath(suffix ...string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get the workspace directory: %w", err)
	}

	fileName := "credentials.csv"
	if len(suffix) > 0 {
		fileName = fmt.Sprintf("credentials-%s.csv", strings.Join(suffix, "-"))
	}
	return filepath.Join(wd, fileName), nil
}

func testAccCheckV3AccessKeyDefaultCredentialFileMode600(suffix ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		path, err := defaultV3AccessKeyCredentialFilePath(suffix...)
		if err != nil {
			return err
		}

		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("credential file not found at %s: %w", path, err)
		}
		if info.Mode().Perm() != 0600 {
			return fmt.Errorf("credential file %s has mode %o, expected 0600", path, info.Mode().Perm())
		}
		return nil
	}
}

func testAccCleanupV3AccessKeyDefaultCredentialFile(suffix ...string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		path, err := defaultV3AccessKeyCredentialFilePath(suffix...)
		if err != nil {
			return err
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("unable to remove credential file %s: %w", path, err)
		}
		return nil
	}
}

func TestAccV3AccessKey_withoutSecretFileInput(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_access_key.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getV3AccessKeyResourceFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			testAccCleanupV3AccessKeyDefaultCredentialFile(name),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccV3AccessKey_withoutSecretFileInput_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|(\.\d{6}Z))$`)),
					resource.TestCheckNoResourceAttr(resourceName, "secret"),
					testAccCheckV3AccessKeyDefaultCredentialFileMode600(name),
				),
			},
		},
	})
}

func testAccV3AccessKey_withoutSecretFileInput_base(name string) string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}

resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = random_string.test.result
  enabled     = true
  description = "Created by terraform script"
}
`, name)
}

func testAccV3AccessKey_withoutSecretFileInput_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_access_key" "test" {
  user_id     = huaweicloud_identity_user.test.id
  description = "Created by terraform script"

  # Clean up the credential file (created by huaweicloud_identity_access_key resource and with a default name) after the
  # test is completed.
  # -f option is used to force the removal of the file (ignoring the error if the file does not exist, which is expected
  # in the acceptance test workflow because the credentials.csv file will be created in current execution directory).
  provisioner "local-exec" {
    command = format("rm -f %%s", abspath("./credentials-${self.user_name}.csv"))
    when    = destroy
  }
}
`, testAccV3AccessKey_withoutSecretFileInput_base(name), name)
}

func TestAccV3AccessKey_withIncorrectSecretFileInput(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_access_key.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getV3AccessKeyResourceFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3AccessKey_withIncorrectSecretFileInput_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|(\.\d{6}Z))$`)),
					resource.TestMatchResourceAttr(resourceName, "secret", regexp.MustCompile(`^[A-Za-z0-9]{40}$`)),
				),
			},
		},
	})
}

func testAccV3AccessKey_withIncorrectSecretFileInput_base(name string) string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}

resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = random_string.test.result
  enabled     = true
  description = "Created by terraform script"
}
`, name)
}

func testAccV3AccessKey_withIncorrectSecretFileInput_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

# Using an invalid storage path will cause the credentials.csv file to fail to generate, but the service will return
# the secret key information through the secret attribute.
resource "huaweicloud_identity_access_key" "test" {
  user_id     = huaweicloud_identity_user.test.id
  description = "Created by terraform script"
  secret_file = "/null/credentials.csv" # Invalid storage path
}
`, testAccV3AccessKey_withIncorrectSecretFileInput_base(name))
}

// generateTestPGPKeyPair creates a PGP key pair for testing and returns the
// base64-encoded public key, private key, and the key fingerprint.
func generateTestPGPKeyPair(t *testing.T) (publicKeyBase64, privateKeyBase64, fingerprint string) {
	t.Helper()

	entity, err := openpgp.NewEntity("TerraformAccTest", "", "terraform-acc-test@example.com", &packet.Config{
		RSABits: 2048,
	})
	if err != nil {
		t.Fatalf("error generating PGP key pair: %s", err)
	}

	// Sign the self-signature on the primary identity and the subkey binding
	// signature before serializing. NewEntity creates the signature structures
	// but does not actually sign them.
	for _, ident := range entity.Identities {
		if err := ident.SelfSignature.SignUserId(ident.Name, entity.PrimaryKey, entity.PrivateKey, &packet.Config{}); err != nil {
			t.Fatalf("error signing self-signature for identity %q: %s", ident.Name, err)
		}
	}
	for _, subkey := range entity.Subkeys {
		if err := subkey.Sig.SignKey(subkey.PublicKey, entity.PrivateKey, &packet.Config{}); err != nil {
			t.Fatalf("error signing subkey binding: %s", err)
		}
	}

	pubBuf := bytes.NewBuffer(nil)
	if err := entity.Serialize(pubBuf); err != nil {
		t.Fatalf("error serializing public key: %s", err)
	}
	publicKeyBase64 = base64.StdEncoding.EncodeToString(pubBuf.Bytes())

	privBuf := bytes.NewBuffer(nil)
	if err := entity.SerializePrivate(privBuf, nil); err != nil {
		t.Fatalf("error serializing private key: %s", err)
	}
	privateKeyBase64 = base64.StdEncoding.EncodeToString(privBuf.Bytes())

	fingerprint = fmt.Sprintf("%x", entity.PrimaryKey.Fingerprint)

	return
}

// decryptEncryptedSecret decrypts the base64-encoded encrypted secret using the
// base64-encoded private key and returns the plaintext.
func decryptEncryptedSecret(encryptedSecret, privateKeyBase64 string) (string, error) {
	privKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 private key: %w", err)
	}

	cryptBytes, err := base64.StdEncoding.DecodeString(encryptedSecret)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 encrypted secret: %w", err)
	}

	entity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(privKeyBytes)))
	if err != nil {
		return "", fmt.Errorf("error parsing private key: %w", err)
	}

	entityList := &openpgp.EntityList{entity}
	md, err := openpgp.ReadMessage(bytes.NewBuffer(cryptBytes), entityList, nil, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting the message: %w", err)
	}

	ptBuf := bytes.NewBuffer(nil)
	if _, err := ptBuf.ReadFrom(md.UnverifiedBody); err != nil {
		return "", fmt.Errorf("error reading the decrypted message: %w", err)
	}

	return ptBuf.String(), nil
}

func testAccCheckV3AccessKeyEncryptedSecretDecryptable(resourceName, privateKeyBase64 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}

		encryptedSecret := rs.Primary.Attributes["encrypted_secret"]
		if encryptedSecret == "" {
			return errors.New("encrypted_secret is empty")
		}

		decrypted, err := decryptEncryptedSecret(encryptedSecret, privateKeyBase64)
		if err != nil {
			return fmt.Errorf("error decrypting encrypted_secret: %s", err)
		}

		if !regexp.MustCompile(`^[A-Za-z0-9]{40}$`).MatchString(decrypted) {
			return fmt.Errorf("decrypted secret does not match expected format, got: %q", decrypted)
		}

		return nil
	}
}

func TestAccV3AccessKey_withPgpKey(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_access_key.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getV3AccessKeyResourceFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	publicKeyBase64, privateKeyBase64, fingerprint := generateTestPGPKeyPair(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			testAccCleanupV3AccessKeyDefaultCredentialFile(name),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccV3AccessKey_withPgpKey_step1(name, publicKeyBase64),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|(\.\d{6}Z))$`)),
					resource.TestCheckResourceAttr(resourceName, "key_fingerprint", fingerprint),
					resource.TestCheckResourceAttrSet(resourceName, "encrypted_secret"),
					resource.TestCheckNoResourceAttr(resourceName, "secret"),
					testAccCheckV3AccessKeyEncryptedSecretDecryptable(resourceName, privateKeyBase64),
				),
			},
		},
	})
}

func testAccV3AccessKey_withPgpKey_base(name string) string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}

resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = random_string.test.result
  enabled     = true
  description = "Created by terraform script"
}
`, name)
}

func testAccV3AccessKey_withPgpKey_step1(name, pgpKey string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_access_key" "test" {
  user_id     = huaweicloud_identity_user.test.id
  description = "Created by terraform script"
  pgp_key     = "%[2]s"

  # Clean up the credential file (created by huaweicloud_identity_access_key resource and with a default name) after the
  # test is completed.
  provisioner "local-exec" {
    command = format("rm -f %%s", abspath("./credentials-${self.user_name}.csv"))
    when    = destroy
  }
}
`, testAccV3AccessKey_withPgpKey_base(name), pgpKey)
}
