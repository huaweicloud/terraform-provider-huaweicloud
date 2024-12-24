package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccProvisionPermissionSet_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_identitycenter_provision_permission_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testProvisionPermissionSet_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "status", "SUCCEEDED"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateIdFunc:       testProvisionPermissionSetImportState(rName),
				ImportStateVerifyIgnore: []string{"permission_set_id", "account_id"},
				ImportStateVerify:       true,
			},
		},
	})
}

func testProvisionPermissionSetImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		if instanceID == "" {
			return "", fmt.Errorf("attribute (instance_id) of resource (%s) not found: %s", name, rs)
		}

		return instanceID + "/" + rs.Primary.ID, nil
	}
}

func testProvisionPermissionSet_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identitycenter_user" "test" {
  identity_store_id = data.huaweicloud_identitycenter_instance.system.identity_store_id
  user_name         = "%[2]s"
  password_mode     = "OTP"
  family_name       = "test_family_name"
  given_name        = "test_given_name"
  display_name      = "test_display_name"
  email             = "email@example.com"
}

resource "huaweicloud_identitycenter_provision_permission_set" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id
  account_id        = huaweicloud_identitycenter_user.test.id
}
`, testPermissionSet_basic(name), name)
}
