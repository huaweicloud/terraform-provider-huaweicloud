package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssUpgradeTargetImages_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_upgrade_target_images.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCssLowEngineVersion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCssUpgradeTargetImages_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "images.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "images.0.engine_type"),
					resource.TestCheckResourceAttrSet(dataSource, "images.0.engine_version"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("engine_type_filter_is_useful", "true"),
					resource.TestCheckOutput("engine_version_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCssUpgradeTargetImages_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_css_upgrade_target_images" "test" {
  cluster_id   = huaweicloud_css_cluster.test.id
  upgrade_type = "cross"
}

locals {
  image_id       = data.huaweicloud_css_upgrade_target_images.test.images[0].id
  engine_type    = data.huaweicloud_css_upgrade_target_images.test.images[0].engine_type
  engine_version = data.huaweicloud_css_upgrade_target_images.test.images[0].engine_version
}

data "huaweicloud_css_upgrade_target_images" "filter_by_id" {
  cluster_id   = huaweicloud_css_cluster.test.id
  upgrade_type = "cross"
  image_id     = local.image_id
}

data "huaweicloud_css_upgrade_target_images" "filter_by_engine_type" {
  cluster_id   = huaweicloud_css_cluster.test.id
  upgrade_type = "cross"
  engine_type  = local.engine_type
}

data "huaweicloud_css_upgrade_target_images" "filter_by_engine_version" {
  cluster_id     = huaweicloud_css_cluster.test.id
  upgrade_type   = "cross"
  engine_version = local.engine_version
}

locals {
  list_by_id             = data.huaweicloud_css_upgrade_target_images.filter_by_id.images
  list_by_engine_type    = data.huaweicloud_css_upgrade_target_images.filter_by_engine_type.images
  list_by_engine_version = data.huaweicloud_css_upgrade_target_images.filter_by_engine_version.images
}

output "id_filter_is_useful" {
  value = length(local.list_by_id) > 0 && alltrue(
    [for v in local.list_by_id[*].id : v == local.image_id]
  )
}

output "engine_type_filter_is_useful" {
  value = length(local.list_by_engine_type) > 0 && alltrue(
    [for v in local.list_by_engine_type[*].engine_type : v == local.engine_type]
  )
}

output "engine_version_filter_is_useful" {
  value = length(local.list_by_engine_version) > 0 && alltrue(
    [for v in local.list_by_engine_version[*].engine_version : v == local.engine_version]
  )
}
`, testAccCssCluster_lowVersion(name))
}

func testAccCssCluster_lowVersion(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "%[3]s"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

}
`, testAccCssBase(name), name, acceptance.HW_CSS_LOW_ENGINE_VERSION)
}
