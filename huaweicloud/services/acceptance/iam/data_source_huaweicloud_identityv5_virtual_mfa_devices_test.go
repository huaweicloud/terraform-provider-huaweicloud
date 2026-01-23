package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccDataV5VirtualMfaDevices_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identityv5_virtual_mfa_devices.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byUserId   = "data.huaweicloud_identityv5_virtual_mfa_devices.filter_by_user_id"
		dcByUserId = acceptance.InitDataSourceCheck(byUserId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5VirtualMfaDevices_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "devices.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'user_id' parameter.
					dcByUserId.CheckResourceExists(),
					resource.TestCheckOutput("is_user_id_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byUserId, "devices.0.enabled",
						"huaweicloud_identityv5_virtual_mfa_device.test", "enabled"),
					resource.TestCheckResourceAttrPair(byUserId, "devices.0.serial_number",
						"huaweicloud_identityv5_virtual_mfa_device.test", "id"),
					resource.TestCheckResourceAttrPair(byUserId, "devices.0.user_id",
						"huaweicloud_identityv5_virtual_mfa_device.test", "user_id"),
				),
			},
		},
	})
}

func testAccDataV5VirtualMfaDevices_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identityv5_virtual_mfa_device" "test" {
  name    = "%[1]s"
  user_id = huaweicloud_identityv5_user.test.id
}

# Without any filter parameters.
data "huaweicloud_identityv5_virtual_mfa_devices" "test" {
  depends_on = [huaweicloud_identityv5_virtual_mfa_device.test]
}

# Filter by 'user_id' parameter.
locals {
  user_id = huaweicloud_identityv5_user.test.id
}

data "huaweicloud_identityv5_virtual_mfa_devices" "filter_by_user_id" {
  user_id = local.user_id

  depends_on = [huaweicloud_identityv5_virtual_mfa_device.test]
}

locals {
  user_id_filter_result = [for v in data.huaweicloud_identityv5_virtual_mfa_devices.filter_by_user_id.devices[*].user_id :
  v == local.user_id]
}

output "is_user_id_filter_useful" {
  value = length(local.user_id_filter_result) > 0 && alltrue(local.user_id_filter_result)
}
`, name)
}
