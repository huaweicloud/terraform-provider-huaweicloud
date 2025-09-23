package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterParameters_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_cluster_parameters.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceClusterParameters_clusterNotFound(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testDataSourceClusterParameters_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "parameters.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "parameters.0.name"),
					resource.TestCheckResourceAttr(dataSource, "parameters.0.values.#", "2"),
					// The values ​​and default values ​​of some parameters are empty.
					resource.TestCheckResourceAttrSet(dataSource, "parameters.0.values.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "parameters.0.readonly"),
					resource.TestCheckResourceAttrSet(dataSource, "parameters.0.restart_required"),
					resource.TestCheckResourceAttrSet(dataSource, "parameters.0.description"),
					// The `parameters.type`, `parameters.unit`, `parameters.value_range`, `parameters.values.value` and
					// `parameters.values.default_value` of some parameters are empty.
				),
			},
		},
	})
}

func testDataSourceClusterParameters_clusterNotFound() string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_parameters" "test" {
  cluster_id = "%s"
}
`, randUUID)
}

func testDataSourceClusterParameters_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_parameters" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
