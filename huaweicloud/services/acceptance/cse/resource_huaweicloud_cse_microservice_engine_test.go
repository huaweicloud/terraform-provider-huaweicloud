package cse

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v2/engines"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getEngineFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CseV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSE V2 client: %s", err)
	}
	return engines.Get(c, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"])
}

func TestAccMicroserviceEngine_basic(t *testing.T) {
	var (
		engine       engines.Engine
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cse_microservice_engine.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&engine,
		getEngineFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMicroserviceEngine_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform test"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_cse_microservice_engine_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "RBAC"),
					resource.TestCheckResourceAttrSet(resourceName, "admin_pass"),
					resource.TestCheckResourceAttr(resourceName, "availability_zones.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_limit"),
					resource.TestCheckResourceAttrSet(resourceName, "service_limit"),
					resource.TestCheckResourceAttrSet(resourceName, "service_registry_addresses.0.private"),
					resource.TestCheckResourceAttrSet(resourceName, "config_center_addresses.0.private"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"admin_pass",
					"extend_params",
				},
			},
		},
	})
}

func testAccMicroserviceEngine_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cse_microservice_engine_flavors" "test" {
  version = "CSE2"
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  
  bandwidth {
    share_type  = "PER"
    size        = 5
    name        = "%[1]s"
    charge_mode = "traffic"
  }
}
`, rName)
}

func testAccMicroserviceEngine_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cse_microservice_engine" "test" {
  name                  = "%[2]s"
  description           = "Created by terraform test"
  flavor                = data.huaweicloud_cse_microservice_engine_flavors.test.flavors[0].id
  network_id            = huaweicloud_vpc_subnet.test.id
  eip_id                = huaweicloud_vpc_eip.test.id
  enterprise_project_id = "0"

  auth_type  = "RBAC"
  admin_pass = "AccTest!123"

  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

}`, testAccMicroserviceEngine_base(rName), rName)
}

func TestAccMicroserviceEngine_withEpsId(t *testing.T) {
	var (
		engine       engines.Engine
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cse_microservice_engine.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&engine,
		getEngineFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMicroserviceEngine_withEpsId(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform test"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_cse_microservice_engine_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "RBAC"),
					resource.TestCheckResourceAttrSet(resourceName, "admin_pass"),
					resource.TestCheckResourceAttr(resourceName, "availability_zones.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(resourceName, "instance_limit"),
					resource.TestCheckResourceAttrSet(resourceName, "service_limit"),
					resource.TestCheckResourceAttrSet(resourceName, "service_registry_addresses.0.private"),
					resource.TestCheckResourceAttrSet(resourceName, "config_center_addresses.0.private"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"admin_pass",
					"extend_params",
				},
				ImportStateIdFunc: testAccEngineResourceImportStateFunc(resourceName),
			},
		},
	})
}

// With enterprise project ID
func testAccEngineResourceImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["enterprise_project_id"] == "" {
			return "", fmt.Errorf("The imported ID specifies an invalid format, want '{id}/{enterprise_project_id}', "+
				"but '%s/%s'", rs.Primary.ID, rs.Primary.Attributes["enterprise_project_id"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.ID, rs.Primary.Attributes["enterprise_project_id"]), nil
	}
}

func testAccMicroserviceEngine_withEpsId(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cse_microservice_engine" "test" {
  name                  = "%[2]s"
  description           = "Created by terraform test"
  flavor                = data.huaweicloud_cse_microservice_engine_flavors.test.flavors[0].id
  network_id            = huaweicloud_vpc_subnet.test.id
  eip_id                = huaweicloud_vpc_eip.test.id
  enterprise_project_id = "%[3]s"

  auth_type  = "RBAC"
  admin_pass = "AccTest!123"

  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

}`, testAccMicroserviceEngine_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccMicroserviceEngine_nacos(t *testing.T) {
	var (
		engine       engines.Engine
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cse_microservice_engine.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&engine,
		getEngineFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMicroserviceEngine_nacos_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "flavor"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "NONE"),
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

func testAccMicroserviceEngine_nacos_step1(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_cse_microservice_engine_flavors" "test" {
  version = "Nacos2"
}

%[1]s

resource "huaweicloud_cse_microservice_engine" "test" {
  name       = "%[2]s"
  # 10 capacity units (500 microservice instances)
  flavor     = format("%%s.10", data.huaweicloud_cse_microservice_engine_flavors.test.flavors[0].id)
  network_id = huaweicloud_vpc_subnet.test.id
  auth_type  = "NONE"
  version    = "Nacos2"
}`, common.TestVpc(rName), rName)
}
