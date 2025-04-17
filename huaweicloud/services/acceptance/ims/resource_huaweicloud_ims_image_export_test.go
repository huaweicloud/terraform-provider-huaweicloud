package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccImageExport_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This resource will export an image file to the OBS bucket,
			// so here we need to set a URL in the format **OBS bucket name:image file name**.
			acceptance.TestAccPreCheckImsImageUrl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// One-time action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccImageExport_basic(),
			},
		},
	})
}

func TestAccImageExport_with_isQuickExport(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This resource will export an image file to the OBS bucket,
			// so here we need to set a URL in the format **OBS bucket name:image file name**.
			acceptance.TestAccPreCheckImsImageUrl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// One-time action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccImageExport_with_isQuickExport(),
			},
		},
	})
}

func testAccImageExport_base(rName string) string {
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
`, common.TestBaseNetwork(rName), rName)
}

func testAccImageExport_basic() string {
	rName := acceptance.RandomAccResourceName()
	rNameSuffix := acctest.RandString(4)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_image_export" "test" {
  image_id    = huaweicloud_ims_ecs_system_image.test.id
  bucket_url  = "%[2]s_%[3]s"
  file_format = "qcow2"
}
`, testAccImageExport_base(rName), acceptance.HW_IMS_IMAGE_URL, rNameSuffix)
}

func testAccImageExport_with_isQuickExport() string {
	rName := acceptance.RandomAccResourceName()
	rNameSuffix := acctest.RandString(4)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_image_export" "test" {
  image_id        = huaweicloud_ims_ecs_system_image.test.id
  bucket_url      = "%[2]s_%[3]s"
  is_quick_export = true
}
`, testAccImageExport_base(rName), acceptance.HW_IMS_IMAGE_URL, rNameSuffix)
}
