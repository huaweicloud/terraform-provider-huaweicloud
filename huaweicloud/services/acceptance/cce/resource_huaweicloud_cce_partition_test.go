package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/partitions"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getPartitionResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CceV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE v1 client: %s", err)
	}
	resp, err := partitions.Get(c, state.Primary.Attributes["cluster_id"],
		state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("unable to find the partition (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccCCEPartition_basic(t *testing.T) {
	var partition partitions.Partitions
	resourceName := "huaweicloud_cce_partition.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	// The availability zone of IES edge partition, example: "cn-south-1-ies-fstxz"
	azName := acceptance.HW_CCE_PARTITION_AZ
	groupName := acceptance.HW_CCE_PARTITION_GROUP

	rc := acceptance.InitResourceCheck(
		resourceName,
		&partition,
		getPartitionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCcePartitionAz(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPartition_basic(randName, azName, groupName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id",
						"huaweicloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", azName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPartitionImportStateIdFunc(azName),
			},
		},
	})
}

func testAccPartitionImportStateIdFunc(azName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_cce_cluster" {
				clusterID = rs.Primary.ID
			}
		}
		if clusterID == "" || azName == "" {
			return "", fmt.Errorf("resource not found: %s/%s", clusterID, azName)
		}
		return fmt.Sprintf("%s/%s", clusterID, azName), nil
	}
}

// testVpc vpc with center and edge zone
func testVpc(name, azName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet_center" {
  name       = "subnet-center"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_vpc_subnet" "subnet_edge" {
  name              = "subnet-edge"
  vpc_id            = huaweicloud_vpc.test.id
  cidr              = "192.168.1.0/24"
  gateway_ip        = "192.168.1.1"
  availability_zone = "%s"
}
`, name, azName)
}

func testAccPartition_Base(randName, azName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                         = "%s"
  cluster_type                 = "VirtualMachine"
  flavor_id                    = "cce.s1.small"
  vpc_id                       = huaweicloud_vpc.test.id
  subnet_id                    = huaweicloud_vpc_subnet.subnet_center.id
  container_network_type       = "eni"
  eni_subnet_id                = huaweicloud_vpc_subnet.subnet_center.ipv4_subnet_id
  eni_subnet_cidr              = huaweicloud_vpc_subnet.subnet_center.cidr
  enable_distribute_management = true
}
`, testVpc(randName, azName), randName)
}

func testAccPartition_basic(randName, azName, groupName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_partition" "test" {
  cluster_id           = huaweicloud_cce_cluster.test.id
  category             = "IES"
  name                 = "%[2]s"
  public_border_group  = "%[3]s"
  partition_subnet_id  = huaweicloud_vpc_subnet.subnet_edge.id
  container_subnet_ids = [huaweicloud_vpc_subnet.subnet_edge.ipv4_subnet_id]
}
`, testAccPartition_Base(randName, azName), azName, groupName)
}
