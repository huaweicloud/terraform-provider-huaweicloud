package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceLogicalClusterRings_basic(t *testing.T) {
	rName := "data.huaweicloud_dws_logical_cluster_rings.test"

	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceLogicalClusterRings_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_rings.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.is_available"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.host_name"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.back_ip"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.cpu_cores"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.memory"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.disk_size"),
				),
			},
		},
	})
}

func testAccDatasourceLogicalClusterRings_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dws_cluster" "test" {
  name              = "%[2]s"
  node_type         = "dwsk2.xlarge"
  number_of_node    = 7
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  user_name         = "test_cluster_admin"
  user_pwd          = "cluster123@!"
}
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourceLogicalClusterRings_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dws_logical_cluster_rings" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id
}
`, testAccDatasourceLogicalClusterRings_base())
}
