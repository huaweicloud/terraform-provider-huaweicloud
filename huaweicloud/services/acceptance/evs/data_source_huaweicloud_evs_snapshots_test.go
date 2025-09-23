package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceEVSsnapshots_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_evs_snapshots.test"
	epsIDTestName := "data.huaweicloud_evs_snapshots.enterprise_project_id_filter"
	serviceTypeTestName := "data.huaweicloud_evs_snapshots.service_type_filter"
	availabilityZoneTestName := "data.huaweicloud_evs_snapshots.availability_zone_filter"
	dedicatedStorageNameTestName := "data.huaweicloud_evs_snapshots.dedicated_storage_name_filter"
	dedicatedStorageIDTestName := "data.huaweicloud_evs_snapshots.dedicated_storage_id_filter"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSnapshots_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "snapshots.0.id"),
					resource.TestCheckResourceAttrSet(rName, "snapshots.0.name"),
					resource.TestCheckResourceAttrSet(rName, "snapshots.0.status"),
					resource.TestCheckResourceAttrSet(rName, "snapshots.0.description"),
					resource.TestCheckResourceAttrSet(rName, "snapshots.0.size"),
					resource.TestCheckResourceAttrSet(rName, "snapshots.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "snapshots.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "snapshots.0.volume_id"),
					resource.TestCheckResourceAttrSet(rName, "snapshots.0.service_type"),
					resource.TestCheckResourceAttrSet(rName, "snapshots.0.progress"),
					resource.TestCheckResourceAttr(rName, "snapshots.0.metadata.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "snapshots.0.metadata.key", "value"),
					resource.TestCheckResourceAttrSet(epsIDTestName, "snapshots.#"),
					resource.TestCheckResourceAttrSet(serviceTypeTestName, "snapshots.#"),
					resource.TestCheckResourceAttrSet(availabilityZoneTestName, "snapshots.#"),
					resource.TestCheckResourceAttrSet(dedicatedStorageNameTestName, "snapshots.#"),
					resource.TestCheckResourceAttrSet(dedicatedStorageIDTestName, "snapshots.#"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),

					resource.TestCheckOutput("volume_id_filter_is_useful", "true"),

					resource.TestCheckOutput("snapshot_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceSnapshots_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_evs_snapshots" "test" {
  depends_on = [huaweicloud_evs_snapshot.test]
}

data "huaweicloud_evs_snapshots" "enterprise_project_id_filter" {
  enterprise_project_id = huaweicloud_evs_volume.test.enterprise_project_id
}

data "huaweicloud_evs_snapshots" "name_filter" {
  name = data.huaweicloud_evs_snapshots.test.snapshots.0.name
}

locals {
  name = data.huaweicloud_evs_snapshots.test.snapshots.0.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_evs_snapshots.name_filter.snapshots) > 0 && alltrue(
    [for v in data.huaweicloud_evs_snapshots.name_filter.snapshots[*].name : v == local.name]
  )  
}

data "huaweicloud_evs_snapshots" "status_filter" {
  status = data.huaweicloud_evs_snapshots.test.snapshots.0.status
}

locals {
  status = data.huaweicloud_evs_snapshots.test.snapshots.0.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_evs_snapshots.status_filter.snapshots) > 0 && alltrue(
    [for v in data.huaweicloud_evs_snapshots.status_filter.snapshots[*].status : v == local.status]
  )  
}

data "huaweicloud_evs_snapshots" "volume_id_filter" {
  volume_id = data.huaweicloud_evs_snapshots.test.snapshots.0.volume_id
}

locals {
  volume_id = data.huaweicloud_evs_snapshots.test.snapshots.0.volume_id
}

output "volume_id_filter_is_useful" {
  value = length(data.huaweicloud_evs_snapshots.volume_id_filter.snapshots) > 0 && alltrue(
    [for v in data.huaweicloud_evs_snapshots.volume_id_filter.snapshots[*].volume_id : v == local.volume_id]
  )  
}

data "huaweicloud_evs_snapshots" "snapshot_id_filter" {
  snapshot_id = data.huaweicloud_evs_snapshots.test.snapshots.0.id
}

locals {
  snapshot_id = data.huaweicloud_evs_snapshots.test.snapshots.0.id
}

output "snapshot_id_filter_is_useful" {
  value = length(data.huaweicloud_evs_snapshots.snapshot_id_filter.snapshots) > 0 && alltrue(
    [for v in data.huaweicloud_evs_snapshots.snapshot_id_filter.snapshots[*].id : v == local.snapshot_id]
  )  
}

data "huaweicloud_evs_snapshots" "service_type_filter" {
  service_type = "evs"
}

data "huaweicloud_evs_snapshots" "availability_zone_filter" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

data "huaweicloud_evs_snapshots" "dedicated_storage_name_filter" {
  dedicated_storage_name = huaweicloud_evs_volume.test.dedicated_storage_name
}

data "huaweicloud_evs_snapshots" "dedicated_storage_id_filter" {
  dedicated_storage_id = huaweicloud_evs_volume.test.dedicated_storage_id
}
`, testAccEvsSnapshotV2_basic(name))
}
