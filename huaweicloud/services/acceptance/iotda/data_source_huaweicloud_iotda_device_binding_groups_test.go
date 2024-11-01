package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDeviceBindingGroups_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_device_binding_groups.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Only the standard and enterprise versions of IoTDA instances support this resource.
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDeviceBindingGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.type"),
				),
			},
		},
	})
}

func testAccDataSourceDeviceBindingGroups_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%s

data "huaweicloud_iotda_device_binding_groups" "test" {
  depends_on = [
    huaweicloud_iotda_device_group.test,
    huaweicloud_iotda_device.test,
  ]

  device_id = huaweicloud_iotda_device.test.id
}
`, testDeviceGroup_basic(name, name))
}
