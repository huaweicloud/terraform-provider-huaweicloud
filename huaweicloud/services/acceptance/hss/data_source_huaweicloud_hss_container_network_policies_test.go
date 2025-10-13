package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerNetworkPolicies_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_container_network_policies.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceContainerNetworkPolicies_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "last_update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
				),
			},
		},
	})
}

const testAccDataSourceContainerNetworkPolicies_basic string = `
data "huaweicloud_hss_container_network_policies" "test" {
  cluster_id = "non-exist-id"
}
`
