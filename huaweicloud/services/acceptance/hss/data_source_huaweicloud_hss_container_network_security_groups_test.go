package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerNetworkSecurityGroups_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_network_security_groups.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceContainerNetworkSecurityGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.0.security_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.0.security_group_name"),
				),
			},
		},
	})
}

func testDataSourceContainerNetworkSecurityGroups_basic() string {
	return `
data "huaweicloud_hss_container_network_security_groups" "test" {
  enterprise_project_id = "0"
}
`
}
