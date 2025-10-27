package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `data_list`.
func TestAccDataSourceClusterProtectPolicies_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_cluster_protect_policies.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceClusterProtectPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "general_policy_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "malicious_image_policy_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_policy_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
				),
			},
		},
	})
}

// The value of `cluster_id` is fake data.
func testAccDataSourceClusterProtectPolicies_basic() string {
	return `
data "huaweicloud_hss_cluster_protect_policies" "test" {
  cluster_id            = "3bd2a82c-4b37-47f3-952b-fa323c22c8e6"
  enterprise_project_id = "0"
}
`
}
