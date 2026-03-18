package cse

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cse"
)

func getMicroserviceEngineFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cse", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSE client: %s", err)
	}

	return cse.GetMicroserviceEngineById(client, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"])
}

func parseInputEnterpriseProjectId(enterpriseProjectId string) string {
	if enterpriseProjectId == "" {
		return "0"
	}

	return enterpriseProjectId
}

func TestAccMicroserviceEngine_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_cse_microservice_engine.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getMicroserviceEngineFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Allow the enterprise project ID to be empty, which will be managed in the default enterprise project
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMicroserviceEngine_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform test"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_cse_microservice_engine_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "RBAC"),
					resource.TestCheckResourceAttrSet(resourceName, "admin_pass"),
					resource.TestCheckResourceAttr(resourceName, "availability_zones.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						parseInputEnterpriseProjectId(acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)),
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
				ImportStateIdFunc: testAccMicroserviceEngineImportStateFunc(resourceName),
			},
		},
	})
}

// If the resource belongs to the default enterprise project, only the resource ID should be passed during import;
// otherwise, both the resource ID and the enterprise project ID should be passed.
func testAccMicroserviceEngineImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["enterprise_project_id"] == "" || rs.Primary.Attributes["enterprise_project_id"] == "0" {
			return rs.Primary.ID, nil
		}
		return fmt.Sprintf("%s/%s", rs.Primary.ID, rs.Primary.Attributes["enterprise_project_id"]), nil
	}
}

func testAccMicroserviceEngine_base(name string) string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 16
  min_numeric      = 1
  min_lower        = 1
  min_upper        = 1
  min_special      = 1
  override_special = "@#"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cse_microservice_engine_flavors" "test" {
  version = "CSE2"
}

%[1]s

resource "huaweicloud_vpc_eip" "test" {
  enterprise_project_id = huaweicloud_vpc.test.enterprise_project_id != "" ? huaweicloud_vpc.test.enterprise_project_id : null

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    size        = 5
    name        = "%[2]s"
    charge_mode = "traffic"
  }
}
`, common.TestVpc(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), name)
}

func testAccMicroserviceEngine_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cse_microservice_engine" "test" {
  name                  = "%[2]s"
  description           = "Created by terraform test"
  flavor                = data.huaweicloud_cse_microservice_engine_flavors.test.flavors[0].id
  network_id            = huaweicloud_vpc_subnet.test.id
  # eip_id                = huaweicloud_vpc_eip.test.id
  availability_zones    = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  enterprise_project_id = huaweicloud_vpc.test.enterprise_project_id != "" ? huaweicloud_vpc.test.enterprise_project_id : null

  auth_type  = "RBAC"
  admin_pass = format("pwdPrefix%%s", random_string.test.result) // Avoid the password starting with a special character
}`, testAccMicroserviceEngine_base(name), name)
}

func TestAccMicroserviceEngine_nacos(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_cse_microservice_engine.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getMicroserviceEngineFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Allow the enterprise project ID to be empty, which will be managed in the default enterprise project
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMicroserviceEngine_nacos_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "flavor"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						parseInputEnterpriseProjectId(acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "NONE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMicroserviceEngineImportStateFunc(resourceName),
			},
		},
	})
}

func testAccMicroserviceEngine_nacos_step1(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cse_microservice_engine_flavors" "test" {
  version = "Nacos2"
}

%[1]s

resource "huaweicloud_cse_microservice_engine" "test" {
  name                  = "%[2]s"
  # 10 capacity units (500 microservice instances)
  flavor                = format("%%s.10", data.huaweicloud_cse_microservice_engine_flavors.test.flavors[0].id)
  network_id            = huaweicloud_vpc_subnet.test.id
  auth_type             = "NONE"
  version               = "Nacos2"
  availability_zones    = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  enterprise_project_id = huaweicloud_vpc.test.enterprise_project_id != "" ? huaweicloud_vpc.test.enterprise_project_id : null
}`, common.TestVpc(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), name)
}
