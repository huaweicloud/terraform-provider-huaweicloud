package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEipAutoProtectStatus_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_eip_auto_protect_status.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a firewall instance ID and the enterprise project ID it belongs to.
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEipAutoProtectStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.available_eip_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.beyond_max_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.eip_protected_self"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.eip_total"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.eip_un_protected"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
				),
			},
		},
	})
}

func testDataSourceEipAutoProtectStatus_base() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = "%s"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

func testDataSourceEipAutoProtectStatus_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cfw_eip_auto_protect_status" "test" {
  object_id             = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  enterprise_project_id = "%[2]s"
}
`, testDataSourceEipAutoProtectStatus_base(), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
