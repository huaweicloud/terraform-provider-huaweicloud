package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLogicalClusters_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_logical_clusters.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			// The cluster_id dose not exist (the cluster ID must be in the standard UUID format).
			{
				Config:      testDataSourceLogicalClusters_expectError(),
				ExpectError: regexp.MustCompile("DWS.0047"),
			},
			// The cluster_id dose not exist (in non-standard UUID format).
			{
				Config:      testDataSourceLogicalClusters_clusterIdNotExist,
				ExpectError: regexp.MustCompile("DWS.0001"),
			},
			{
				Config: testDataSourceLogicalClusters_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "logical_clusters.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "logical_clusters.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "logical_clusters.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "logical_clusters.0.status"),
					resource.TestCheckOutput("is_existed_logical_cluster", "true"),
					resource.TestCheckOutput("assert_ring_hosts", "true"),
				),
			},
		},
	})
}

func testDataSourceLogicalClusters_expectError() string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_logical_clusters" "test" {
  cluster_id = "%s"
}
`, randUUID)
}

const testDataSourceLogicalClusters_clusterIdNotExist = `
data "huaweicloud_dws_logical_clusters" "test" {
  cluster_id = "not-found-cluster-id"
}
`

func testDataSourceLogicalClusters_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dws_logical_cluster_rings" "test" {
  cluster_id = "%[1]s"
}

resource "huaweicloud_dws_logical_cluster" "test" {
  cluster_id           = "%[1]s"
  logical_cluster_name = "%[2]s"

  cluster_rings {
    dynamic "ring_hosts" {
      for_each = local.ring_hosts
      content {
        host_name = ring_hosts.value.host_name
        back_ip   = ring_hosts.value.back_ip
        cpu_cores = ring_hosts.value.cpu_cores
        memory    = ring_hosts.value.memory
        disk_size = ring_hosts.value.disk_size
      }
    }
  }
}

data "huaweicloud_dws_logical_clusters" "test" {
  depends_on = [
    huaweicloud_dws_logical_cluster.test
  ]

  cluster_id = "%[1]s"
}

locals {
  ring_hosts = data.huaweicloud_dws_logical_cluster_rings.test.cluster_rings.0.ring_hosts[*]
  logical_cluster_id = huaweicloud_dws_logical_cluster.test.id
  logical_clusters      = data.huaweicloud_dws_logical_clusters.test.logical_clusters
  logical_cluster_ids   = local.logical_clusters[*].id
  logical_cluster_names = local.logical_clusters[*].name
  ring_hosts_result     = [for v in local.logical_clusters : length(v.cluster_rings[0].ring_hosts) == length(local.ring_hosts) if v.id
  == local.logical_cluster_id]
}

# Assert that the query results contain the current logical cluster ID and name.
output "is_existed_logical_cluster" {
  value = contains(local.logical_cluster_ids, local.logical_cluster_id) && contains(local.logical_cluster_names, "%[2]s")
}

# Assert ring hosts of the current resource.
output "assert_ring_hosts" {
  value = length(local.ring_hosts_result) == 1 && alltrue(local.ring_hosts_result)
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}
