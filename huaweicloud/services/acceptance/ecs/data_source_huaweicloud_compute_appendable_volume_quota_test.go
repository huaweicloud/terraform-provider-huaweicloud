package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEcsComputeAppendableVolumeQuota_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_appendable_volume_quota.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeAppendableVolumeQuota_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quota_count"),
					resource.TestCheckResourceAttrSet(dataSource, "free_scsi"),
					resource.TestCheckResourceAttrSet(dataSource, "free_blk"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeAppendableVolumeQuota_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}
`, testAccCompute_data, rName)
}

func testDataSourceEcsComputeAppendableVolumeQuota_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_compute_appendable_volume_quota" "test" {
  server_id = huaweicloud_compute_instance.test.id
}
`, testDataSourceEcsComputeAppendableVolumeQuota_base(name))
}
