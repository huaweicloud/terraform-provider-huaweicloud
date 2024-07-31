package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccEvsVolumesDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_evs_volumes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolumesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "volumes.#", "5"),

					resource.TestCheckOutput("is_volume_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_availability_zone_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccEvsVolumesDataSource_base(rName string) string {
	return fmt.Sprintf(`
variable "volume_configuration" {
  type = list(object({
    suffix      = string
    size        = number
    device_type = string
    multiattach = bool
  }))
  default = [
    {suffix = "vbd_normal_volume", size = 100, device_type = "VBD", multiattach = false},
    {suffix = "vbd_share_volume", size = 100, device_type = "VBD", multiattach = true},
    {suffix = "scsi_normal_volume", size = 100, device_type = "SCSI", multiattach = false},
    {suffix = "scsi_share_volume", size = 100, device_type = "SCSI", multiattach = true},
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

func testAccEvsVolumesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_evs_volumes" "test" {
  depends_on = [huaweicloud_compute_volume_attach.test]

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  server_id         = huaweicloud_compute_instance.test.id
  status            = "in-use"
}

# Filter using volume ID
data "huaweicloud_evs_volumes" "volume_id_filter" {
  volume_id = data.huaweicloud_evs_volumes.test.volumes.0.id
}

locals {
  volume_id = data.huaweicloud_evs_volumes.test.volumes.0.id
}

output "is_volume_id_filter_useful" {
  value = length(data.huaweicloud_evs_volumes.volume_id_filter.volumes) > 0 && alltrue(
    [for v in data.huaweicloud_evs_volumes.volume_id_filter.volumes[*].id : v == local.volume_id]
  )
}

# Filter using name
data "huaweicloud_evs_volumes" "name_filter" {
  name = data.huaweicloud_evs_volumes.test.volumes.0.name
}

# The name parameter is a fuzzy query, so only the number of query results is verified
output "is_name_filter_useful" {
  value = length(data.huaweicloud_evs_volumes.name_filter.volumes) > 0
}

# Filter using availability_zone
data "huaweicloud_evs_volumes" "availability_zone_filter" {
  availability_zone = data.huaweicloud_evs_volumes.test.volumes.0.availability_zone
}

locals {
  availability_zone = data.huaweicloud_evs_volumes.test.volumes.0.availability_zone
}

output "is_availability_zone_filter_useful" {
  value = length(data.huaweicloud_evs_volumes.availability_zone_filter.volumes) > 0 && alltrue(
    [for v in data.huaweicloud_evs_volumes.availability_zone_filter.volumes[*].availability_zone : v == local.availability_zone]
  )
}

# Filter using enterprise_project_id
data "huaweicloud_evs_volumes" "enterprise_project_id_filter" {
  enterprise_project_id = data.huaweicloud_evs_volumes.test.volumes.0.enterprise_project_id
}

locals {
  enterprise_project_id = data.huaweicloud_evs_volumes.test.volumes.0.enterprise_project_id
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_evs_volumes.enterprise_project_id_filter.volumes) > 0 && alltrue(
    [for v in data.huaweicloud_evs_volumes.enterprise_project_id_filter.volumes[*].enterprise_project_id : v == local.enterprise_project_id]
  )
}
`, testAccEvsVolumesDataSource_base(rName))
}
