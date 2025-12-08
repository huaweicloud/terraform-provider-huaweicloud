package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccIMSV21ImageExport_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIMSV21ImageExport_basic(),
			},
		},
	})
}

func TestAccIMSV21ImageExport_with_isQuickExport(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIMSV21ImageExport_with_isQuickExport(),
			},
		},
	})
}

func testAccIMSV21ImageExport_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_name         = "Ubuntu 18.04 server 64bit"
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_ims_ecs_system_image" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_compute_instance.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[2]s"
  acl           = "private"
  force_destroy = true
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccIMSV21ImageExport_basic() string {
	rName := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_imsv21_image_export" "test" {
  image_id    = huaweicloud_ims_ecs_system_image.test.id
  bucket_url  = "${huaweicloud_obs_bucket.test.bucket}:%[2]s.qcow2"
  file_format = "qcow2"
}
`, testAccIMSV21ImageExport_base(rName), rName)
}

func testAccIMSV21ImageExport_with_isQuickExport() string {
	rName := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_imsv21_image_export" "test" {
  image_id        = huaweicloud_ims_ecs_system_image.test.id
  bucket_url      = "${huaweicloud_obs_bucket.test.bucket}:%[2]s.zvhd2"
  is_quick_export = true
}
`, testAccIMSV21ImageExport_base(rName), rName)
}
