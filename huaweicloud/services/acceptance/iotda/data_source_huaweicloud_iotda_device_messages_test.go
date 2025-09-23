package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDeviceMessages_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_device_messages.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDeviceMessages_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "messages.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "messages.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "messages.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "messages.0.message"),
					resource.TestCheckResourceAttrSet(dataSourceName, "messages.0.encoding"),
					resource.TestCheckResourceAttrSet(dataSourceName, "messages.0.payload_format"),
					resource.TestCheckResourceAttrSet(dataSourceName, "messages.0.topic"),
					resource.TestCheckResourceAttr(dataSourceName, "messages.0.properties.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "messages.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "messages.0.created_time"),
				),
			},
		},
	})
}

func testAccDataSourceDeviceMessages_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%s

data "huaweicloud_iotda_device_messages" "test" {
  depends_on = [
    huaweicloud_iotda_device.test,
    huaweicloud_iotda_device_message.test,
  ]

  device_id = huaweicloud_iotda_device.test.id
}
`, testDeviceMessage_basic(name))
}
