package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDbWdrSnapshotCollectionResults_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
			acceptance.TestAccPreCheckGaussDBTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbWdrSnapshotCollectionResults_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.#"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.file_size"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.wdr_type"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.job_create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.download_url"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.notes"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.file_path"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.obs_bucket.#"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.obs_bucket.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.obs_bucket.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.obs_bucket.0.url"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.obs_bucket.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "wdr_snapshots.0.obs_bucket.0.domain_id"),
					resource.TestCheckOutput("start_time_filter_is_useful", "true"),
					resource.TestCheckOutput("end_time_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("wdr_type_filter_is_useful", "true"),
					resource.TestCheckOutput("job_start_time_filter_is_useful", "true"),
					resource.TestCheckOutput("job_end_time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbWdrSnapshotCollectionResults_basic() string {
	return fmt.Sprintf(`
%[2]s

data "huaweicloud_gaussdb_wdr_snapshot_collection_results" "test" {
  depends_on = [huaweicloud_gaussdb_wdr_snapshot_collect.test]

  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_wdr_snapshot_collection_results" "start_time_filter" {
  depends_on = [huaweicloud_gaussdb_wdr_snapshot_collect.test]

  instance_id = "%[1]s"
  start_time  = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].start_time
}
locals {
  start_time = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].start_time
}
output "start_time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_wdr_snapshot_collection_results.start_time_filter.wdr_snapshots) > 0
}

data "huaweicloud_gaussdb_wdr_snapshot_collection_results" "end_time_filter" {
  depends_on = [huaweicloud_gaussdb_wdr_snapshot_collect.test]

  instance_id = "%[1]s"
  end_time    = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].end_time
}
locals {
  end_time = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].end_time
}
output "end_time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_wdr_snapshot_collection_results.end_time_filter.wdr_snapshots) > 0
}

data "huaweicloud_gaussdb_wdr_snapshot_collection_results" "status_filter" {
  depends_on = [huaweicloud_gaussdb_wdr_snapshot_collect.test]

  instance_id = "%[1]s"
  status      = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].status
}
locals {
  status = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].status
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_wdr_snapshot_collection_results.status_filter.wdr_snapshots) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_wdr_snapshot_collection_results.status_filter.wdr_snapshots[*].status :
  v == local.status]
  )
}

data "huaweicloud_gaussdb_wdr_snapshot_collection_results" "wdr_type_filter" {
  depends_on = [huaweicloud_gaussdb_wdr_snapshot_collect.test]

  instance_id = "%[1]s"
  wdr_type    = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].wdr_type
}
locals {
  wdr_type = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].wdr_type
}
output "wdr_type_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_wdr_snapshot_collection_results.wdr_type_filter.wdr_snapshots) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_wdr_snapshot_collection_results.wdr_type_filter.wdr_snapshots[*].wdr_type :
  v == local.wdr_type]
  )
}

data "huaweicloud_gaussdb_wdr_snapshot_collection_results" "job_start_time_filter" {
  depends_on = [huaweicloud_gaussdb_wdr_snapshot_collect.test]

  instance_id    = "%[1]s"
  job_start_time = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].job_create_time
}
locals {
  job_start_time = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].job_create_time
}
output "job_start_time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_wdr_snapshot_collection_results.job_start_time_filter.wdr_snapshots) > 0
}

data "huaweicloud_gaussdb_wdr_snapshot_collection_results" "job_end_time_filter" {
  depends_on = [huaweicloud_gaussdb_wdr_snapshot_collect.test]

  instance_id  = "%[1]s"
  job_end_time = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].job_create_time
}
locals {
  job_end_time = data.huaweicloud_gaussdb_wdr_snapshot_collection_results.test.wdr_snapshots[0].job_create_time
}
output "job_end_time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_wdr_snapshot_collection_results.job_end_time_filter.wdr_snapshots) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, testAccGaussDBWdrSnapshotCollect_basic(), testAccGaussDBWdrSnapshotCollect_nodes())
}
