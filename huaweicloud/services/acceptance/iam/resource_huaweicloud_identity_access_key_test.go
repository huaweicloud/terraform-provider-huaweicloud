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

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/credentials"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/encryption"
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

func TestAccV3AccessKey_withPgpKey(t *testing.T) {
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
			acceptance.TestAccPreCheckPGPKeys(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3AccessKey_withPgpKey_step1(name, acceptance.HW_PGP_PUBLIC_KEY),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestMatchResourceAttr(resourceName, "secret", regexp.MustCompile(`^[A-Za-z0-9]{40}$`)),
					resource.TestMatchResourceAttr(resourceName, "encrypted_secret", regexp.MustCompile(`^[A-Za-z0-9+/=]+$`)),
					testAccCheckV3AccessKeyPgpEncryption(resourceName),
				),
			},
		},
	})
}

func testAccV3AccessKey_withPgpKey_step1(name, publicKey string) string {
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

# Using an invalid storage path ensures the secret is stored in tfstate for verification.
resource "huaweicloud_identity_access_key" "test" {
  user_id     = huaweicloud_identity_user.test.id
  description = "Created by terraform script"
  secret_file = "/null/credentials.csv"
  pgp_key     = "%[2]s"
}
`, name, publicKey)
}

func testAccCheckV3AccessKeyPgpEncryption(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fingerprint, err := encryption.GetPGPFingerprint(acceptance.HW_PGP_PUBLIC_KEY)
		if err != nil {
			return fmt.Errorf("error computing PGP key fingerprint: %w", err)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		gotFingerprint := rs.Primary.Attributes["key_fingerprint"]
		if gotFingerprint != fingerprint {
			return fmt.Errorf("The fingerprint of the PGP key used to encrypt the secret is incorrect, want '%v', but got '%v'", fingerprint, gotFingerprint)
		}

		encryptedSecret := rs.Primary.Attributes["encrypted_secret"]
		if encryptedSecret == "" {
			return errors.New("encrypted_secret is empty")
		}

		secret := rs.Primary.Attributes["secret"]
		if secret == "" {
			return errors.New("secret is empty")
		}

		decrypted, err := decryptPGPMessage(encryptedSecret, acceptance.HW_PGP_PRIVATE_KEY)
		if err != nil {
			return fmt.Errorf("error decrypting encrypted_secret: %w", err)
		}
		if decrypted != secret {
			return fmt.Errorf("The decrypted secret is incorrect, want '%v', but got '%v'", secret, decrypted)
		}

		return nil
	}
}

func decryptPGPMessage(encryptedSecret, privateKey string) (string, error) {
	privKeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error decoding private key: %w", err)
	}
	cryptBytes, err := base64.StdEncoding.DecodeString(encryptedSecret)
	if err != nil {
		return "", fmt.Errorf("error decoding encrypted secret: %w", err)
	}

	entity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(privKeyBytes)))
	if err != nil {
		return "", fmt.Errorf("error parsing private key: %w", err)
	}

	md, err := openpgp.ReadMessage(bytes.NewBuffer(cryptBytes), &openpgp.EntityList{entity}, nil, nil)
	if err != nil {
		return "", fmt.Errorf("error reading encrypted message: %w", err)
	}

	ptBuf := bytes.NewBuffer(nil)
	if _, err = ptBuf.ReadFrom(md.UnverifiedBody); err != nil {
		return "", fmt.Errorf("error reading decrypted message body: %w", err)
	}

	return ptBuf.String(), nil
}
