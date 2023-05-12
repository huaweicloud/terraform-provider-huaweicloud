package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceFirewalls_basic(t *testing.T) {
	rName := "data.huaweicloud_cfw_firewalls.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceFirewalls_basic(),
				Check: resource.ComposeTestCheckFunc(
					// only check whether the API can be called successfully,
					// more attributes check will be added
					// when the resource to create a firewall is available
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "records.0.fw_instance_id", acceptance.HW_CFW_INSTANCE_ID),
				),
			},
			{
				Config: testAccDatasourceFirewalls_empty(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccDatasourceFirewalls_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = "%s"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

func testAccDatasourceFirewalls_empty() string {
	return `
data "huaweicloud_cfw_firewalls" "test" {}
`
}
