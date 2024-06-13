package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIamIdentityVirtualMfaDevices_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_identity_virtual_mfa_devices.basic"
	dataSource2 := "data.huaweicloud_identity_virtual_mfa_devices.filter_by_user_id"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceIamIdentityVirtualMfaDevices_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_user_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceIamIdentityVirtualMfaDevices_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_virtual_mfa_devices" "basic" {
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
  value = length(data.huaweicloud_identity_virtual_mfa_devices.basic.virtual_mfa_devices) > 0
}

output "is_user_id_filter_useful" {
  value = alltrue(local.user_id_filter_result) && length(local.user_id_filter_result) > 0
}
`, testAccIdentityVirtualMFADevice_basic(name), acceptance.HW_USER_ID)
}
