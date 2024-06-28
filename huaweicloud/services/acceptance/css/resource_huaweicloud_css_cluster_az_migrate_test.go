package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCssClusterAzMigrate_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCssAzMigrateAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCssClusterAzMigrate_basic(rName),
			},
		},
	})
}

func testAccCssClusterAzMigrate_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  az0 = data.huaweicloud_availability_zones.test.names[0]
  az1 = data.huaweicloud_availability_zones.test.names[1]
  az2 = data.huaweicloud_availability_zones.test.names[2]
}

resource "huaweicloud_css_cluster_az_migrate" "test" {
  cluster_id           = huaweicloud_css_cluster.test.id
  instance_type        = "ess"
  source_az            = local.az0
  target_az            = "${local.az0},${local.az1},${local.az2}"
  migrate_type         = "multi_az_change"
  agency               = "%[2]s"
  indices_backup_check = true
}
`, testAccCssClusterAz(rName), acceptance.HW_CSS_AZ_MIGRATE_AGENCY)
}

func testAccCssClusterAz(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 3
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  lifecycle {
    ignore_changes = [
	  availability_zone
    ]
  }
}
`, testAccCssBase(name), name)
}
