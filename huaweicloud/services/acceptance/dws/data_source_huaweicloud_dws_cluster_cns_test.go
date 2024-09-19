package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcClusterCns_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_cluster_cns.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	uuid, _ := uuid.GenerateUUID()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceClusterCns_basic(uuid),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testDataSourceClusterCns_basic(acceptance.HW_DWS_CLUSTER_ID),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "max_num"),
					resource.TestCheckResourceAttrSet(dataSource, "min_num"),
					resource.TestCheckResourceAttrSet(dataSource, "cns.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cns.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "cns.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "cns.0.private_ip"),
				),
			},
		},
	})
}

func testDataSourceClusterCns_basic(clusterId string) string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_cns" "test" {
  cluster_id = "%s"
}
`, clusterId)
}
