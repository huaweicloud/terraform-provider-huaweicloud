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
		return nil, fmtp.Errorf("error creating WAF dedicated client: %s", err)
	}
	return instances.GetWithEpsId(client, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"])
}

func TestAccDedicatedInstance_basic(t *testing.T) {
	var (
		instance     instances.DedicatedInstance
		resourceName = "huaweicloud_waf_dedicated_instance.test"
		name         = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getWafDedicatedInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			// Configure two enterprise projects to test enterprise project migration.
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cpu_architecture", "x86"),
					resource.TestCheckResourceAttr(resourceName, "specification_code", "waf.instance.professional"),
					resource.TestCheckResourceAttr(resourceName, "security_group.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "run_status", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "upgradable", "0"),
					resource.TestCheckResourceAttr(resourceName, "res_tenant", "true"),
					resource.TestCheckResourceAttr(resourceName, "anti_affinity", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "server_id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ecs_flavor"),
					resource.TestCheckResourceAttrSet(resourceName, "available_zone"),
				),
			},
			{
				Config: testAccWafDedicatedInstance_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", name+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "cpu_architecture", "x86"),
					resource.TestCheckResourceAttr(resourceName, "specification_code", "waf.instance.professional"),
					resource.TestCheckResourceAttr(resourceName, "security_group.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "run_status", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "upgradable", "0"),
					resource.TestCheckResourceAttr(resourceName, "res_tenant", "true"),
					resource.TestCheckResourceAttr(resourceName, "anti_affinity", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "server_id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ecs_flavor"),
					resource.TestCheckResourceAttrSet(resourceName, "available_zone"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testWAFResourceImportState(resourceName),
				ImportStateVerifyIgnore: []string{"res_tenant", "anti_affinity", "tags"},
			},
		},
	})
}

func testAccWafDedicatedInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_dedicated_instance" "test" {
  name                  = "%[2]s"
  available_zone        = data.huaweicloud_availability_zones.test.names[1]
  specification_code    = "waf.instance.professional"
  ecs_flavor            = data.huaweicloud_compute_flavors.test.ids[0]
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  res_tenant            = true
  anti_affinity         = true
  enterprise_project_id = "%[3]s"
  
  tags = {
    foo = "bar"
    key = "value"
  }

  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}
`, common.TestBaseComputeResources(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafDedicatedInstance_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_dedicated_instance" "test" {
  name                  = "%[2]s_updated"
  available_zone        = data.huaweicloud_availability_zones.test.names[1]
  specification_code    = "waf.instance.professional"
  ecs_flavor            = data.huaweicloud_compute_flavors.test.ids[0]
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  res_tenant            = true
  anti_affinity         = true
  enterprise_project_id = "%[3]s"
  
  tags = {
    foo = "bar"
    key = "value"
  }

  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}
`, common.TestBaseComputeResources(name), name, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}
