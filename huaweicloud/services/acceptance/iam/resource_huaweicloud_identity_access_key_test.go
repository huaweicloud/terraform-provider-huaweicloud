package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/credentials"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getAccessKeyResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	iamClient, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	found, err := credentials.Get(iamClient, state.Primary.ID).Extract()
	if err != nil {
		return nil, err
	}

	if found.AccessKey != state.Primary.ID {
		return nil, fmt.Errorf("Access Key not found")
	}
	return found, nil
}

func TestAccAccessKey_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_access_key.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getAccessKeyResourceFunc)

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
				Config: testAccAccessKey_basic_step1(name),
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
				Config: testAccAccessKey_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "inactive"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|(\.\d{6}Z))$`)),
				),
			},
		},
	})
}

func testAccAccessKey_basic_base(name string) string {
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

func testAccAccessKey_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_access_key" "test" {
  user_id     = huaweicloud_identity_user.test.id
  description = "Created by terraform script"
  secret_file = abspath("./credentials.csv")
}
`, testAccAccessKey_basic_base(name))
}

func testAccAccessKey_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_access_key" "test" {
  user_id     = huaweicloud_identity_user.test.id	
  description = "Updated by terraform script"
  secret_file = abspath("./credentials.csv")
  status      = "inactive"

  # Clean up the credentials.csv file (created by huaweicloud_identity_access_key resource) after the test is completed.
  provisioner "local-exec" {
    command = "rm credentials.csv"
    when    = destroy
  }
}
`, testAccAccessKey_basic_base(name))
}

func TestAccAccessKey_withoutSecretFileInput(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_access_key.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getAccessKeyResourceFunc)

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
				Config: testAccAccessKey_withoutSecretFileInput_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|(\.\d{6}Z))$`)),
					resource.TestCheckNoResourceAttr(resourceName, "secret"),
				),
			},
		},
	})
}

func testAccAccessKey_withoutSecretFileInput_base(name string) string {
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

func testAccAccessKey_withoutSecretFileInput_step1(name string) string {
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
`, testAccAccessKey_withoutSecretFileInput_base(name), name)
}

func TestAccAccessKey_withIncorrectSecretFileInput(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_access_key.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getAccessKeyResourceFunc)

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
				Config: testAccAccessKey_withIncorrectSecretFileInput_step1(name),
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

func testAccAccessKey_withIncorrectSecretFileInput_base(name string) string {
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

func testAccAccessKey_withIncorrectSecretFileInput_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

# Using an invalid storage path will cause the credentials.csv file to fail to generate, but the service will return
# the secret key information through the secret attribute.
resource "huaweicloud_identity_access_key" "test" {
  user_id     = huaweicloud_identity_user.test.id
  description = "Created by terraform script"
  secret_file = "/null/credentials.csv" # Invalid storage path
}
`, testAccAccessKey_withIncorrectSecretFileInput_base(name))
}
