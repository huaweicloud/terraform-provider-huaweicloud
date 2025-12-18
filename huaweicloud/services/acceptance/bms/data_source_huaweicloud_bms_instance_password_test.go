package bms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBmsInstancePassword_basic(t *testing.T) {
	dataSource := "data.huaweicloud_bms_instance_password.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckBmsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBmsInstancePassword_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "password"),
				),
			},
		},
	})
}

func testDataSourceBmsInstancePassword_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_bms_instance_password" "test" {
  server_id = "%s"
}
`, acceptance.HW_BMS_INSTANCE_ID)
}
