package cfw

import (
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
				),
			},
		},
	})
}

func testAccDatasourceFirewalls_basic() string {
	return `
data "huaweicloud_cfw_firewalls" "test" {
  service_type = "0"
}
`
}
