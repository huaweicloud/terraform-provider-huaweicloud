package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5VirtualMfaDevices_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_virtual_mfa_devices.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5VirtualMfaDevices_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.serial_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.user_id"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5VirtualMfaDevices_basic() string {
	return `
data "huaweicloud_identityv5_virtual_mfa_devices" "test" {}
`
}
