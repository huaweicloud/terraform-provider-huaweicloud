package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	instances "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func getWafDedicatedInstanceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.WafDedicatedV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating HuaweiCloud WAF dedicated client : %s", err)
	}
	return instances.GetInstance(client, state.Primary.ID)
}

func TestAccWafDedicatedInstance_basic(t *testing.T) {
	var instance instances.DedicatedInstance
	resourceName := "huaweicloud_waf_dedicated_instance.instance_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getWafDedicatedInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstanceV1_conf(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cpu_architecture", "x86"),
					resource.TestCheckResourceAttr(resourceName, "specification_code", "waf.instance.professional"),
					resource.TestCheckResourceAttr(resourceName, "security_group.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "run_status", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "upgradable", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "server_id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ecs_flavor"),
					resource.TestCheckResourceAttrSet(resourceName, "available_zone"),
				),
			},
			{
				Config: testAccWafDedicatedInstance_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "cpu_architecture", "x86"),
					resource.TestCheckResourceAttr(resourceName, "specification_code", "waf.instance.professional"),
					resource.TestCheckResourceAttr(resourceName, "security_group.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "run_status", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "upgradable", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "server_id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ecs_flavor"),
					resource.TestCheckResourceAttrSet(resourceName, "available_zone"),
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

func TestAccWafDedicatedInstance_elb_model(t *testing.T) {
	var instance instances.DedicatedInstance
	resourceName := "huaweicloud_waf_dedicated_instance.instance_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getWafDedicatedInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstance_elb_model(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cpu_architecture", "x86"),
					resource.TestCheckResourceAttr(resourceName, "specification_code", "waf.instance.professional"),
					resource.TestCheckResourceAttr(resourceName, "security_group.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "run_status", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "upgradable", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "server_id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ecs_flavor"),
					resource.TestCheckResourceAttrSet(resourceName, "available_zone"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "group_id",
						"${huaweicloud_waf_instance_group.group_1.id}"),
				),
			},
		},
	})
}

func baseDependResource(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "%s_waf"
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_vpc_subnet" "vpc_subnet_1" {
  name       = "%s_waf"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.vpc_1.id
}

resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "%s_waf"
  description = "terraform security group acceptance test"
}

data "huaweicloud_availability_zones" "zones" {}

data "huaweicloud_compute_flavors" "flavors" {
  availability_zone = data.huaweicloud_availability_zones.zones.names[1]
  performance_type  = "normal"
  cpu_core_count    = 2
}
`, name, name, name)
}

func testAccWafDedicatedInstanceV1_conf(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%s"
  available_zone     = data.huaweicloud_availability_zones.zones.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.flavors.ids[0]
  vpc_id             = huaweicloud_vpc.vpc_1.id
  subnet_id          = huaweicloud_vpc_subnet.vpc_subnet_1.id
  
  security_group = [
    huaweicloud_networking_secgroup.secgroup.id
  ]
}
`, baseDependResource(name), name)
}

func testAccWafDedicatedInstance_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%s_updated"
  available_zone     = data.huaweicloud_availability_zones.zones.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.flavors.ids[0]
  vpc_id             = huaweicloud_vpc.vpc_1.id
  subnet_id          = huaweicloud_vpc_subnet.vpc_subnet_1.id
  
  security_group = [
    huaweicloud_networking_secgroup.secgroup.id
  ]
}
`, baseDependResource(name), name)
}

func testAccWafDedicatedInstance_elb_model(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_instance_group" "group_1" {
  name   = "%s"
  vpc_id = huaweicloud_vpc.vpc_1.id
}

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%s"
  available_zone     = data.huaweicloud_availability_zones.zones.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.flavors.ids[0]
  vpc_id             = huaweicloud_vpc.vpc_1.id
  subnet_id          = huaweicloud_vpc_subnet.vpc_subnet_1.id
  group_id           = huaweicloud_waf_instance_group.group_1.id
  
  security_group = [
    huaweicloud_networking_secgroup.secgroup.id
  ]
}
`, baseDependResource(name), name, name)
}
