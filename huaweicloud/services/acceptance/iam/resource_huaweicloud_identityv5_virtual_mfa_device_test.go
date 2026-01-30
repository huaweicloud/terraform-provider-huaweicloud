package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getV5VirtualMFADeviceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM Client: %s", err)
	}

	return iam.GetV5VirtualMfaDevice(client, state.Primary.Attributes["user_id"])
}

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5VirtualMFADevice_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_identityv5_virtual_mfa_device.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5VirtualMFADeviceResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV5VirtualMFADevice_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "user_id", "huaweicloud_identityv5_user.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "base32_string_seed"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV5VirtualMFADeviceImportStateFunc(rName),
				ImportStateVerifyIgnore: []string{
					"name", "base32_string_seed",
				},
			},
		},
	})
}

func testAccV5VirtualMFADevice_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identityv5_virtual_mfa_device" "test" {
  name    = "%[1]s"
  user_id = huaweicloud_identityv5_user.test.id
}
`, name)
}

func testAccV5VirtualMFADeviceImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		userID := rs.Primary.Attributes["user_id"]
		if userID == "" {
			return "", fmt.Errorf("attribute (user_id) of Resource (%s) not found", rName)
		}

		return userID, nil
	}
}
