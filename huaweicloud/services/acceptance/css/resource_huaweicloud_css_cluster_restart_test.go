package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCssClusterRestart_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster_restart.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssClusterRestart_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func TestAccCssClusterRestart_rolling_restart(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster_restart.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssClusterRestart_rolling_restart(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccCssClusterRestart_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster_restart" "test" {
  cluster_id = huaweicloud_css_cluster.test.id
  type       = "role"
  value      = "ess"
}
`, testAccCssCluster_basic(rName, "Test@passw0rd", 7, "bar"))
}

func testAccCssClusterRestart_rolling_restart(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster_restart" "test" {
  cluster_id = huaweicloud_css_cluster.test.id
  type       = "role"
  value      = "ess"
  is_rolling = true
}
`, testAccCssClusterRestart_build(rName))
}

func testAccCssClusterRestart_build(rName string) string {
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
