package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	instances "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func getWafDedicatedInstanceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.WafDedicatedV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating HuaweiCloud WAF dedicated client : %s", err)
	}
	return instances.GetWithEpsId(client, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"])
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

func TestAccWafDedicatedInstance_withEpsId(t *testing.T) {
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
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstance_epsId(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(
						resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testAccWafDedicatedInstance_epsId(name, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(
						resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFResourceImportState(resourceName),
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

func TestAccWafDedicatedInstance_resourceTenant(t *testing.T) {
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
				Config: testAccWafDedicatedInstanceV1_resourceTenant(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cpu_architecture", "x86"),
					resource.TestCheckResourceAttr(resourceName, "specification_code", "waf.instance.professional"),
					resource.TestCheckResourceAttr(resourceName, "security_group.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "run_status", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "upgradable", "0"),
					resource.TestCheckResourceAttr(resourceName, "res_tenant", "true"),
					resource.TestCheckResourceAttr(resourceName, "anti_affinity", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "server_id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "available_zone"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"res_tenant", "anti_affinity"},
			},
		},
	})
}

func testAccWafDedicatedInstanceV1_conf(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%s"
  available_zone     = data.huaweicloud_availability_zones.test.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.test.ids[0]
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  
  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}
`, common.TestBaseComputeResources(name), name)
}

func testAccWafDedicatedInstance_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%s_updated"
  available_zone     = data.huaweicloud_availability_zones.test.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.test.ids[0]
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  
  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}
`, common.TestBaseComputeResources(name), name)
}

func testAccWafDedicatedInstance_epsId(name, epsId string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name                  = "%s"
  available_zone        = data.huaweicloud_availability_zones.test.names[1]
  specification_code    = "waf.instance.professional"
  ecs_flavor            = data.huaweicloud_compute_flavors.test.ids[0]
  enterprise_project_id = "%s"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  
  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}
`, common.TestBaseComputeResources(name), name, epsId)
}

func testAccWafDedicatedInstance_elb_model(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_instance_group" "group_1" {
  name   = "%[2]s"
  vpc_id = huaweicloud_vpc.test.id
}

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%[2]s"
  available_zone     = data.huaweicloud_availability_zones.test.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.test.ids[0]
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  group_id           = huaweicloud_waf_instance_group.group_1.id
  
  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}
`, common.TestBaseComputeResources(name), name)
}

func testAccWafDedicatedInstanceV1_resourceTenant(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%s"
  available_zone     = data.huaweicloud_availability_zones.test.names[1]
  specification_code = "waf.instance.professional"
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  res_tenant         = true
  anti_affinity      = true

  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}
`, common.TestBaseComputeResources(name), name)
}
