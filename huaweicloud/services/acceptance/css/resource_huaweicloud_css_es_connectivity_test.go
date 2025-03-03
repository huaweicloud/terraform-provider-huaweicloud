package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEsConnectivity_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEsConnectivity_basic(name),
			},
		},
	})
}

func testAccEsConnectivity_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  count          = 2
  name           = "%[2]s_${count.index}"
  engine_version = "7.10.2"
  password       = "Test@passw0rd"

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

resource "huaweicloud_css_es_connectivity" "test" {
  source_cluster_id = huaweicloud_css_cluster.test[0].id
  target_cluster_id = huaweicloud_css_cluster.test[1].id
}
`, testAccCssBase(name), name)
}
