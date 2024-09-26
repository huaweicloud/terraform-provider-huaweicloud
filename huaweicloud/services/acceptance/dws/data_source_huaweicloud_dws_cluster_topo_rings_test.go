package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterTopoRings_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_cluster_topo_rings.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceClusterTopoRings_clusterNotExist(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testDataSourceClusterTopoRings_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "rings.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_exist_instance", "true"),
					resource.TestCheckResourceAttrSet(dataSource, "rings.0.instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rings.0.instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "rings.0.instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "rings.0.instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "rings.0.instances.0.manage_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "rings.0.instances.0.internal_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "rings.0.instances.0.traffic_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "rings.0.instances.0.availability_zone"),
				),
			},
		},
	})
}

func testDataSourceClusterTopoRings_clusterNotExist() string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_topo_rings" "test" {
  cluster_id = "%s"
}
`, randUUID)
}

func testDataSourceClusterTopoRings_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_topo_rings" "test" {
  cluster_id = "%s"
}

// For a DWS cluster, there are at least three nodes.
output "is_exist_instance" {
  value = length(data.huaweicloud_dws_cluster_topo_rings.test.rings[0].instances) >= 3
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
