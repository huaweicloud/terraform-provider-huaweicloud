package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusters_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_clusters.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceClusters_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "clusters.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_exist_cluster_id", "true"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.node_type"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.number_of_node"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.security_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.tags.%"),
					resource.TestMatchResourceAttr(dataSource, "clusters.0.endpoints.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(dataSource, "clusters.0.nodes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.sub_status"),
					resource.TestMatchResourceAttr(dataSource, "clusters.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "clusters.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testDataSourceClusters_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_clusters" "test" {}

locals {
  cluster_ids = data.huaweicloud_dws_clusters.test.clusters[*].id
}

output "is_exist_cluster_id" {
  value = contains(local.cluster_ids, "%s")
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
