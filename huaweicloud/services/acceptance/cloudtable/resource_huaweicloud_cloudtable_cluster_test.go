package cloudtable

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cloudtable/v2/clusters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getClusterResourceObj(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CloudtableV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud CloudTable v2 client: %s", err)
	}
	return clusters.Get(c, state.Primary.ID)
}

func TestAccCloudtableCluster_basic(t *testing.T) {
	var cluster clusters.Cluster
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cloudtable_cluster.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&cluster,
		getClusterResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudtableCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", acceptance.HW_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "storage_type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "storage_size", "10240"),
					resource.TestCheckResourceAttr(resourceName, "rs_num", "4"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"availability_zone", "network_id",
				},
			},
		},
	})
}

func testAccCloudtableCluster_base(rName string) string {
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()

	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "%s"
  gateway_ip = "%s"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}`, rName, randCidr, rName, randCidr, randGatewayIp, rName)
}

func testAccCloudtableCluster_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cloudtable_cluster" "test" {
  availability_zone = "%s"
  name              = "%s"
  storage_type      = "ULTRAHIGH"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  rs_num            = 4
}
`, testAccCloudtableCluster_base(rName), acceptance.HW_AVAILABILITY_ZONE, rName)
}
