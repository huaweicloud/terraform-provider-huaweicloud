package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to the lack of the `risk_id` parameter, this test case was not successfully executed.
func TestAccDataSourceContainerClusterRiskAffectedResources_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_cluster_risk_affected_resources.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSClusterRiskId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceContainerClusterRiskAffectedResources_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.hit_rule"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.first_scan_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.last_scan_time"),
				),
			},
		},
	})
}

func testAccDataSourceContainerClusterRiskAffectedResources_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_container_cluster_risk_affected_resources" "test" {
  risk_id = "%s"
}
`, acceptance.HW_HSS_CLUSTER_RISK_ID)
}
