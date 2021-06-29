package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/identity/v3.0/credentials"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIdentityAccessKey_basic(t *testing.T) {
	var cred credentials.Credential
	var userName = fmt.Sprintf("acc-user-%s", acctest.RandString(5))
	resourceName := "huaweicloud_identity_access_key.key_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAccessKey_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityAccessKeyExists(resourceName, &cred),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description", "access key by terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: testAccIdentityAccessKey_update(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "inactive"),
					resource.TestCheckResourceAttr(resourceName, "description", "access key by terraform updated"),
				),
			},
		},
	})
}

func testAccCheckIdentityAccessKeyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	iamClient, err := config.IAMV3Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_access_key" {
			continue
		}

		_, err := credentials.Get(iamClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Access Key still exists")
		}
	}

	return nil
}

func testAccCheckIdentityAccessKeyExists(n string, cred *credentials.Credential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		iamClient, err := config.IAMV3Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
		}

		found, err := credentials.Get(iamClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.AccessKey != rs.Primary.ID {
			return fmtp.Errorf("Access Key not found")
		}

		*cred = *found

		return nil
	}
}

func testAccIdentityAccessKey_basic(userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "user_1" {
  name        = "%s"
  password    = "password123@!"
  enabled     = true
  description = "tested by terraform"
}

resource "huaweicloud_identity_access_key" "key_1" {
  user_id     = huaweicloud_identity_user.user_1.id
  description = "access key by terraform"
  secret_file = "./credentials.csv"
}
`, userName)
}

func testAccIdentityAccessKey_update(userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "user_1" {
  name        = "%s"
  password    = "password123@!"
  enabled     = true
  description = "tested by terraform"
}

resource "huaweicloud_identity_access_key" "key_1" {
  user_id     = huaweicloud_identity_user.user_1.id
  description = "access key by terraform updated"
  secret_file = "./credentials.csv"
  status      = "inactive"
}
`, userName)
}
