package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLogicalClusterVolumes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_logical_cluster_volumes.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsLogicalModeClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			// The DWS cluster ID does not exist.
			{
				Config:      testDataSourceLogicalClusterVolumes_expectError(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testDataSourceLogicalClusterVolumes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "volumes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "volumes.0.logical_cluster_name"),
					resource.TestCheckResourceAttrSet(dataSource, "volumes.0.percentage"),
					resource.TestCheckResourceAttrSet(dataSource, "volumes.0.usage"),
					resource.TestCheckResourceAttrSet(dataSource, "volumes.0.total"),
				),
			},
		},
	})
}

func testDataSourceLogicalClusterVolumes_expectError() string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_logical_cluster_volumes" "test" {
  cluster_id = "%s"
}
`, randUUID)
}

func testDataSourceLogicalClusterVolumes_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_logical_cluster_volumes" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_DWS_LOGICAL_MODE_CLUSTER_ID)
}
