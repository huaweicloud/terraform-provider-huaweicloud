package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityVirtualMfaDevices_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identity_virtual_mfa_devices.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byUserId   = "data.huaweicloud_identity_virtual_mfa_devices.filter_by_user_id"
		dcByUserId = acceptance.InitDataSourceCheck(byUserId)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityVirtualMfaDevices_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),

					// filter by user_id
					dcByUserId.CheckResourceExists(),
					resource.TestCheckOutput("is_user_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccIdentityVirtualMfaDevices_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_virtual_mfa_devices" "test" {
  depends_on = [
    huaweicloud_identity_virtual_mfa_device.test
  ]
}

data "huaweicloud_identity_virtual_mfa_devices" "filter_by_user_id" {
  user_id = "%[2]s"

  depends_on = [
    huaweicloud_identity_virtual_mfa_device.test
  ]
}

locals {
  user_id_filter_result = [
    for v in data.huaweicloud_identity_virtual_mfa_devices.filter_by_user_id.virtual_mfa_devices[*].user_id : v == "%[2]s"
  ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_identity_virtual_mfa_devices.test.virtual_mfa_devices) > 0
}

output "is_user_id_filter_useful" {
  value = alltrue(local.user_id_filter_result) && length(local.user_id_filter_result) > 0
}
`, testAccIdentityVirtualMFADevice_basic(name), acceptance.HW_USER_ID)
}
