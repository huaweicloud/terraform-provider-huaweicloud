package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscDeviceMaskAlgorithms_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_device_mask_algorithms.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscDeviceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscDeviceMaskAlgorithms_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "mask_algorithms.#"),
				),
			},
		},
	})
}

func testAccDataSourceDscDeviceMaskAlgorithms_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_device_mask_algorithms" "test" {
  device_id = "%[1]s"
}
`, acceptance.HW_DSC_DEVICE_ID)
}
