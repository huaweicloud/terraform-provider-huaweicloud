package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	instances "github.com/huaweicloud/golangsdk/openstack/waf_hw/v1/premium_instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWafDedicatedInstanceV1_basic(t *testing.T) {
	var instance instances.DedicatedInstance
	resourceName := "huaweicloud_waf_dedicated_instance.instance_1"
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckWafDedicatedInstanceV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstanceV1_conf(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedInstanceV1Exists(resourceName, &instance),
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
				Config: testAccWafDedicatedInstanceV1_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedInstanceV1Exists(resourceName, &instance),
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

func testAccCheckWafDedicatedInstanceV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	c, err := config.WafDedicatedV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud WAF dedicated client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_waf_dedicated_instance" {
			continue
		}
		_, err := instances.GetInstance(c, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Waf domain still exists")
		}
	}
	return nil
}

func testAccCheckWafDedicatedInstanceV1Exists(n string, instance *instances.DedicatedInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		c, err := config.WafDedicatedV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating HuaweiCloud WAF dedicated client: %s", err)
		}

		found, err := instances.GetInstance(c, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Waf dedicated instance not found")
		}
		*instance = *found
		return nil
	}
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

func testAccWafDedicatedInstanceV1_update(name string) string {
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
