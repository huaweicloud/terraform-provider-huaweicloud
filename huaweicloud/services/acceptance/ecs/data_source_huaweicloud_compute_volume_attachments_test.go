package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEcsComputeVolumeAttachments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_volume_attachments.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeVolumeAttachments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.#"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.device"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.server_id"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.volume_id"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeVolumeAttachments_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_images" "test" {
  flavor_id = data.huaweicloud_compute_flavors.test.ids[0]

  os         = "Ubuntu"
  visibility = "public"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[1]s"
  image_id           = data.huaweicloud_images_images.test.images[0].id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_evs_volume" "test" {
  count = 2
  name              = "%[1]s_${count.index}"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 10
}

resource "huaweicloud_compute_volume_attach" "test" {
  count = 2

  instance_id           = huaweicloud_compute_instance.test.id
  volume_id             = huaweicloud_evs_volume.test[count.index].id
  delete_on_termination = "true"
}
`, name)
}

func testDataSourceEcsComputeVolumeAttachments_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_compute_volume_attachments" "test" {
  depends_on = [
    huaweicloud_compute_volume_attach.test[0],
    huaweicloud_compute_volume_attach.test[1]
  ]
  server_id = huaweicloud_compute_instance.test.id
}
`, testDataSourceEcsComputeVolumeAttachments_base(name))
}
