package cdm

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cdm/v1/clusters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getCdmClusterResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.CdmV11Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating CDM v1 client, err=%s", err)
	}
	return clusters.Get(client, state.Primary.ID)
}

func TestAccResourceCdmCluster_basic(t *testing.T) {
	var obj clusters.ClusterCreateOpts
	resourceName := "huaweicloud_cdm_cluster.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdmClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdmCluster_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttr(resourceName, "status", "Normal"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCdmCluster_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cdm_flavors" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  vpc_id     = huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "%s"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_cdm_cluster" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  flavor_id         = data.huaweicloud_cdm_flavors.test.flavors[0].id
  name              = "%s"
  security_group_id = huaweicloud_networking_secgroup.secgroup.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
}
`, name, name, name, name)
}

func TestAccResourceCdmCluster_all(t *testing.T) {
	var obj clusters.ClusterCreateOpts
	resourceName := "huaweicloud_cdm_cluster.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdmClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdmCluster_all(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "is_auto_off", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "Normal"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
				),
			},
			{
				Config: testAccCdmCluster_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "is_auto_off", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "Normal"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"email", "phone_num"},
			},
		},
	})
}

func testAccCdmCluster_all(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cdm_flavors" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  vpc_id     = huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "%s"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_cdm_cluster" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  flavor_id         = data.huaweicloud_cdm_flavors.test.flavors[0].id
  name              = "%s"
  security_group_id = huaweicloud_networking_secgroup.secgroup.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  is_auto_off       = true
  email             = ["test@test.com"]
  phone_num         = ["12345678910"]
}
`, name, name, name, name)
}

func testAccCdmCluster_update(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cdm_flavors" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  vpc_id     = huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "%s"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_cdm_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  flavor_id          = data.huaweicloud_cdm_flavors.test.flavors[0].id
  name               = "%s"
  security_group_id  = huaweicloud_networking_secgroup.secgroup.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  email              = ["test@test.com"]
  phone_num          = ["12345678910"]
  schedule_boot_time = "00:00:00"
  schedule_off_time  = "10:00:00"
}
`, name, name, name, name)
}
