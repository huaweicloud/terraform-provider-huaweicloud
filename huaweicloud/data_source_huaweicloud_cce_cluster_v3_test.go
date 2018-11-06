package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCCEClusterV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCCE(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCCEClusterV3DataSource_cluster,
			},
			resource.TestStep{
				Config: testAccCCEClusterV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3DataSourceID("data.huaweicloud_cce_cluster_v3.clusters"),
					resource.TestCheckResourceAttr("data.huaweicloud_cce_cluster_v3.clusters", "name", "huaweicloud-cce"),
					resource.TestCheckResourceAttr("data.huaweicloud_cce_cluster_v3.clusters", "status", "Available"),
					resource.TestCheckResourceAttr("data.huaweicloud_cce_cluster_v3.clusters", "cluster_type", "VirtualMachine"),
				),
			},
		},
	})
}

func testAccCheckCCEClusterV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cluster data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("cluster data source ID not set ")
		}

		return nil
	}
}

var testAccCCEClusterV3DataSource_cluster = fmt.Sprintf(`
resource "huaweicloud_cce_cluster_v3" "cluster_1" {
  name = "huaweicloud-cce"
  cluster_type = "VirtualMachine"
  flavor = "cce.s1.small"
  cluster_version = "v1.7.3-r10"
  vpc_id = "%s"
  subnet_id = "%s"
  container_network_type = "overlay_l2"
}`, OS_VPC_ID, OS_SUBNET_ID)

var testAccCCEClusterV3DataSource_basic = fmt.Sprintf(`
%s
data "huaweicloud_cce_cluster_v3" "clusters" {
  name = "${huaweicloud_cce_cluster_v3.cluster_1.name}"
}
`, testAccCCEClusterV3DataSource_cluster)
