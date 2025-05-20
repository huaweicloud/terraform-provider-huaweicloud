package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceV3Volumes_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_evsv3_volumes.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		rName          = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV3Volumes_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.links.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.attachments.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.bootable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.volume_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.metadata.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.multiattach"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.iops.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.throughput.#"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_sort_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_metadata_filter_useful", "true"),
					resource.TestCheckOutput("is_availability_zone_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceV3Volumes_base(rName string) string {
	return fmt.Sprintf(`
variable "volume_configuration" {
  type = list(object({
    suffix      = string
    size        = number
    device_type = string
    multiattach = bool
  }))
  default = [
    { suffix = "vbd_normal_volume", size = 100, device_type = "VBD", multiattach = false },
    { suffix = "vbd_share_volume", size = 100, device_type = "VBD", multiattach = true },
    { suffix = "scsi_normal_volume", size = 100, device_type = "SCSI", multiattach = false },
    { suffix = "scsi_share_volume", size = 100, device_type = "SCSI", multiattach = true },
  ]
}

%[1]s

resource "huaweicloud_compute_instance" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[2]s"
  image_id          = data.huaweicloud_images_image.test.id
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]

  system_disk_type = "SSD"
  system_disk_size = 50

  security_group_ids = [
    huaweicloud_networking_secgroup.test.id
  ]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_evs_volume" "test" {
  count = length(var.volume_configuration)

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SSD"
  name              = "%[2]s_${var.volume_configuration[count.index].suffix}"
  size              = var.volume_configuration[count.index].size
  device_type       = var.volume_configuration[count.index].device_type
  multiattach       = var.volume_configuration[count.index].multiattach

  tags = {
    index = tostring(count.index)
  }
}

resource "huaweicloud_compute_volume_attach" "test" {
  count = length(huaweicloud_evs_volume.test)

  instance_id = huaweicloud_compute_instance.test.id
  volume_id   = huaweicloud_evs_volume.test[count.index].id
}
`, common.TestBaseComputeResources(rName), rName)
}
func testAccDataSourceV3Volumes_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_evsv3_volumes" "test" {
  depends_on = [huaweicloud_compute_volume_attach.test]
}

# Filter using name.
data "huaweicloud_evsv3_volumes" "name_filter" {
  name = data.huaweicloud_evsv3_volumes.test.volumes.0.name
}

# The name parameter is a fuzzy query, so only the number of query results is verified.
output "is_name_filter_useful" {
  value = length(data.huaweicloud_evsv3_volumes.name_filter.volumes) > 0
}

# Filter using sort_key and sort_dir.
locals {
  sort_key = "created_at"
}

data "huaweicloud_evsv3_volumes" "sort_asc_filter" {
  depends_on = [huaweicloud_compute_volume_attach.test]

  sort_key = local.sort_key
  sort_dir = "asc"
}

data "huaweicloud_evsv3_volumes" "sort_desc_filter" {
  depends_on = [huaweicloud_compute_volume_attach.test]

  sort_key = local.sort_key
  sort_dir = "desc"
}

locals {
  asc_first_id = data.huaweicloud_evsv3_volumes.sort_asc_filter.volumes[0].id
  desc_length  = length(data.huaweicloud_evsv3_volumes.sort_desc_filter.volumes)
  desc_last_id = data.huaweicloud_evsv3_volumes.sort_desc_filter.volumes[local.desc_length - 1].id
}

output "is_sort_filter_useful" {
  value = local.desc_length > 0 && local.asc_first_id == local.desc_last_id
}

# Filter using status.
data "huaweicloud_evsv3_volumes" "status_filter" {
  status = data.huaweicloud_evsv3_volumes.test.volumes.0.status
}

locals {
  status = data.huaweicloud_evsv3_volumes.test.volumes.0.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_evsv3_volumes.status_filter.volumes) > 0 && alltrue(
    [for v in data.huaweicloud_evsv3_volumes.status_filter.volumes[*].status : v == local.status]
  )
}

# Filter using metadata.
locals {
  metadata = urlencode("{\"hw:passthrough\": \"true\"}")
}

data "huaweicloud_evsv3_volumes" "metadata_filter" {
  depends_on = [huaweicloud_compute_volume_attach.test]

  metadata = local.metadata
}

output "is_metadata_filter_useful" {
  value = length(data.huaweicloud_evsv3_volumes.metadata_filter.volumes) > 0 && alltrue(
    [for v in data.huaweicloud_evsv3_volumes.metadata_filter.volumes[*].metadata : lookup(v, "hw:passthrough", "false") == "true"]
  )
}

# Filter using availability_zone.
data "huaweicloud_evsv3_volumes" "availability_zone_filter" {
  availability_zone = data.huaweicloud_evsv3_volumes.test.volumes.0.availability_zone
}

locals {
  availability_zone = data.huaweicloud_evsv3_volumes.test.volumes.0.availability_zone
}

output "is_availability_zone_filter_useful" {
  value = length(data.huaweicloud_evsv3_volumes.availability_zone_filter.volumes) > 0 && alltrue(
    [for v in data.huaweicloud_evsv3_volumes.availability_zone_filter.volumes[*].availability_zone : v == local.availability_zone]
  )
}

# Filter using non existent name.
data "huaweicloud_evsv3_volumes" "not_found" {
  name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_evsv3_volumes.not_found.volumes) == 0
}
`, testAccDataSourceV3Volumes_base(rName))
}
