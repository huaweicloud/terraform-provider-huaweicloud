package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this acceptance test, make sure that the LTS log function has been enabled.
func TestAccDataSourceClusterLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_cluster_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceClusterLogs_clusterNotFound(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testDataSourceClusterLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "logs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.access_url"),
				),
			},
		},
	})
}

func testDataSourceClusterLogs_clusterNotFound() string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_logs" "test" {
  cluster_id = "%s"
}
`, randUUID)
}

func testDataSourceClusterLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_logs" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
