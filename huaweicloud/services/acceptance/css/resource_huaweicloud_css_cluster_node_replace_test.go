package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCssClusterNodeReplace_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCssReplaceAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCssClusterNodeReplace_basic(rName),
			},
		},
	})
}

func testAccCssClusterNodeReplace_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster_node_replace" "test" {
  cluster_id   = huaweicloud_css_cluster.test.id
  node_id      = huaweicloud_css_cluster.test.nodes[0].id
  agency       = "%[2]s"
  migrate_data = true
}
`, testAccCssClusterNodeReplace_clusterBase(rName), acceptance.HW_CSS_REPLACE_AGENCY)
}

func testAccCssClusterNodeReplace_clusterBase(rName string) string {
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
}
`, testAccCssBase(rName), rName)
}
