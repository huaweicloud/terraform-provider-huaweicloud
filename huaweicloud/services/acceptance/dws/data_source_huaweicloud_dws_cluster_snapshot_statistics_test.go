package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDwsClusterSnapshotStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_cluster_snapshot_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceClusterSnapshotStatistics_clusterNotFound(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testDataSourceClusterSnapshotStatistics_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_total_capacity_set", "true"),
					resource.TestCheckOutput("is_used_capacity_set", "true"),
				),
			},
		},
	})
}

func testDataSourceClusterSnapshotStatistics_clusterNotFound() string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_snapshot_statistics" "test" {
  cluster_id = "%s"
}
`, randUUID)
}

func testDataSourceClusterSnapshotStatistics_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_snapshot" "test" {
  cluster_id = "%[1]s"
  name       = "%[2]s"
}

data "huaweicloud_dws_cluster_snapshot_statistics" "test" {
  depends_on = [huaweicloud_dws_snapshot.test]
  cluster_id = "%[1]s"
}

locals {
  statistics    = data.huaweicloud_dws_cluster_snapshot_statistics.test.statistics
  free_capaticy = [for v in local.statistics : v if v.name == "storage.free"]
  used_capaticy = [for v in local.statistics : v if v.name == "storage.used"]
}

output "is_total_capacity_set" {
  value = try(local.free_capaticy[0].value > 0 && local.free_capaticy[0].unit != "", false)
}

output "is_used_capacity_set" {
  value = try(local.used_capaticy[0].value > 0 && local.used_capaticy[0].unit != "", false)
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}
