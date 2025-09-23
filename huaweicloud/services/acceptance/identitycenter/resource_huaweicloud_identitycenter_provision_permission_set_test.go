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
			acceptance.TestAccPreCheckIdentityCenterAccountId(t)
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
				ResourceName:      rName,
				ImportState:       true,
				ImportStateIdFunc: testProvisionPermissionSetImportState(rName),
				ImportStateVerify: true,
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

resource "huaweicloud_identitycenter_provision_permission_set" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.test.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id
  account_id        = "%[2]s"

  depends_on = [huaweicloud_identitycenter_account_assignment.test]
}
`, testAccountAssignment_basic(name), acceptance.HW_IDENTITY_CENTER_ACCOUNT_ID)
}
