package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRgcDriftDetails_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_drift_details.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRgcDriftDetails_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSource, "drift_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "drift_details.0.drift_message"),
					resource.TestCheckResourceAttrSet(dataSource, "drift_details.0.drift_target_id"),
					resource.TestCheckResourceAttrSet(dataSource, "drift_details.0.drift_target_type"),
					resource.TestCheckResourceAttrSet(dataSource, "drift_details.0.drift_type"),
					resource.TestCheckResourceAttrSet(dataSource, "drift_details.0.managed_account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "drift_details.0.parent_organizational_unit_id"),
					resource.TestCheckResourceAttrSet(dataSource, "drift_details.0.solve_solution"),
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

const testAccDataSourceRgcDriftDetails_basic = `
data "huaweicloud_rgc_drift_details" "test" {
}
`
