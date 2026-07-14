package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceBigdataInstance_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dsc_bigdata_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBigdataInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "asset_name", name),
					resource.TestCheckResourceAttr(rName, "ds_type", "Elasticsearch"),
					resource.TestCheckResourceAttr(rName, "ins_type", "ECS"),
					resource.TestCheckResourceAttr(rName, "scan_metadata", "false"),
				),
			},
		},
	})
}

func testAccResourceBigdataInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_bigdata_instance" "test" {
  depends_on = [
    huaweicloud_compute_instance.test,
  ]

  asset_name        = "%[2]s"
  ds_type           = "Elasticsearch"
  ds_name           = "1"
  ds_address        = huaweicloud_compute_instance.test.access_ip_v4
  ds_port           = 8080
  ds_version        = 5
  ins_type          = "ECS"
  ins_id            = huaweicloud_compute_instance.test.id
  ins_name          = huaweicloud_compute_instance.test.name
  ds_user           = "admin"
  ds_password       = "dsc@123"
  vpc_id            = data.huaweicloud_vpc_subnet.test.vpc_id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  scan_metadata     = false
}
`, testAccResourceBigdataInstance_base(name), name)
}

func testAccResourceBigdataInstance_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[1]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  system_disk_type = "SAS"
  system_disk_size = 50
}
`, name)
}
