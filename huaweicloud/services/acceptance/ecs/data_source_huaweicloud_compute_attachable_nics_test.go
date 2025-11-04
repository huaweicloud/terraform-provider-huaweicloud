package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceEcsComputeAttachableNics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_attachable_nics.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeAttachableNics_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "attachable_quantity.#"),
					resource.TestCheckResourceAttrSet(dataSource, "attachable_quantity.0.free_efi_nic"),
					resource.TestCheckResourceAttrSet(dataSource, "attachable_quantity.0.free_scsi"),
					resource.TestCheckResourceAttrSet(dataSource, "attachable_quantity.0.free_blk"),
					resource.TestCheckResourceAttrSet(dataSource, "attachable_quantity.0.free_disk"),
					resource.TestCheckResourceAttrSet(dataSource, "attachable_quantity.0.free_nic"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.#"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.port_state"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.fixed_ips.#"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.fixed_ips.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.fixed_ips.0.ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.net_id"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.port_id"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.mac_addr"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.delete_on_termination"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.preserve_on_delete"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.min_rate"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.multiqueue_num"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeAttachableNics_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

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
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_images.test.images[0].id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceEcsComputeAttachableNics_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_compute_attachable_nics" "test" {
  server_id = huaweicloud_compute_instance.test.id
}
`, testDataSourceEcsComputeAttachableNics_base(name))
}
