package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRecycleBinVolumesDetail_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_evs_recycle_bin_volumes_detail.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case needs to ensure that there is a volume in the EVS recycle bin.
			acceptance.TestAccPreCheckEVSRecycleBinEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecycleBinVolumesDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.multiattach"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.bootable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.service_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.volume_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.enterprise_project_id"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_availability_zone_filter_useful", "true"),
					resource.TestCheckOutput("is_service_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceRecycleBinVolumesDetail_basic() string {
	return `
data "huaweicloud_evs_recycle_bin_volumes_detail" "test" {}

# Filter using name.
data "huaweicloud_evs_recycle_bin_volumes_detail" "name_filter" {
  name = data.huaweicloud_evs_recycle_bin_volumes_detail.test.volumes.0.name
}

locals {
  name = data.huaweicloud_evs_recycle_bin_volumes_detail.test.volumes.0.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_evs_recycle_bin_volumes_detail.name_filter.volumes) > 0 && alltrue(
    [for v in data.huaweicloud_evs_recycle_bin_volumes_detail.name_filter.volumes[*].name : v == local.name]
  )
}

# Filter using status.
data "huaweicloud_evs_recycle_bin_volumes_detail" "status_filter" {
  status = data.huaweicloud_evs_recycle_bin_volumes_detail.test.volumes.0.status
}

locals {
  status = data.huaweicloud_evs_recycle_bin_volumes_detail.test.volumes.0.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_evs_recycle_bin_volumes_detail.status_filter.volumes) > 0 && alltrue(
    [for v in data.huaweicloud_evs_recycle_bin_volumes_detail.status_filter.volumes[*].status : v == local.status]
  )
}

# Filter using availability_zone.
data "huaweicloud_evs_recycle_bin_volumes_detail" "availability_zone_filter" {
  availability_zone = data.huaweicloud_evs_recycle_bin_volumes_detail.test.volumes.0.availability_zone
}

locals {
  availability_zone = data.huaweicloud_evs_recycle_bin_volumes_detail.test.volumes.0.availability_zone
}

output "is_availability_zone_filter_useful" {
  value = length(data.huaweicloud_evs_recycle_bin_volumes_detail.availability_zone_filter.volumes) > 0 && alltrue(
    [for v in data.huaweicloud_evs_recycle_bin_volumes_detail.availability_zone_filter.volumes[*].availability_zone : v == local.availability_zone]
  )
}

# Filter using service_type.
data "huaweicloud_evs_recycle_bin_volumes_detail" "service_type_filter" {
  service_type = data.huaweicloud_evs_recycle_bin_volumes_detail.test.volumes.0.service_type
}

locals {
  service_type = data.huaweicloud_evs_recycle_bin_volumes_detail.test.volumes.0.service_type
}

output "is_service_type_filter_useful" {
  value = length(data.huaweicloud_evs_recycle_bin_volumes_detail.service_type_filter.volumes) > 0 && alltrue(
    [for v in data.huaweicloud_evs_recycle_bin_volumes_detail.service_type_filter.volumes[*].service_type : v == local.service_type]
  )
}
`
}
