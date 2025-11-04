package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsPolicyStatesStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_policy_states_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRmsPolicyStatesStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "value.#"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.total_resource_count"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.non_compliant_resource_count"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.total_policy_count"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.non_compliant_policy_count"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.statistic_date"),
				),
			},
		},
	})
}

func testDataSourceRmsPolicyStatesStatistics_basic() string {
	return `
data "huaweicloud_rms_policy_states_statistics" "test" {}
`
}
