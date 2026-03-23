package cse

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cse"
)

func getMicroserviceEngineConfigurationFunc(_ *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		// Querying microservice engine configurations requires building a client based on the microservice engine's connection address,
		// which does not use IAM authentication.
		client   = common.NewCustomClient(true, state.Primary.Attributes["connect_address"])
		authInfo = cse.MicroserviceEngineAuthInfo{
			AuthAddress:         getAuthAddress(state.Primary.Attributes),
			AdminUser:           state.Primary.Attributes["admin_user"],
			AdminPass:           state.Primary.Attributes["admin_pass"],
			EnterpriseProjectId: state.Primary.Attributes["enterprise_project_id"],
		}
	)
	return cse.GetMicroserviceEngineConfiguration(client, authInfo, state.Primary.ID)
}

// Before testing, please bind the EIP and open the access rules according to the resource document appendix.
func TestAccMicroserviceEngineConfiguration_basic(t *testing.T) {
	var (
		configuration interface{}

		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cse_microservice_engine_configuration.test"
		rc           = acceptance.InitResourceCheck(resourceName, &configuration, getMicroserviceEngineConfigurationFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSEMicroserviceEngineID(t)
			acceptance.TestAccPreCheckCSEMicroserviceEngineAdminPassword(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMicroserviceEngineConfiguration_basic_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key", randName),
					resource.TestCheckResourceAttr(resourceName, "value_type", "json"),
					resource.TestCheckResourceAttr(resourceName, "value", "{\"foo\":\"bar\"}"),
					resource.TestCheckResourceAttr(resourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", getAcceptanceEpsId()),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccMicroserviceEngineConfiguration_basic_step2(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "value", "{\"foo\":\"baar\"}"),
					resource.TestCheckResourceAttr(resourceName, "status", "disabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMicroserviceEngineConfigurationImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccMicroserviceEngineConfigurationImportStateIdFunc(resName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var authAddr, connAddr, keyName, username, password, enterpriseProjectId string

		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", resName)
		}

		authAddr = rs.Primary.Attributes["auth_address"]
		connAddr = rs.Primary.Attributes["connect_address"]
		keyName = rs.Primary.Attributes["key"]
		username = rs.Primary.Attributes["admin_user"]
		password = rs.Primary.Attributes["admin_pass"]
		enterpriseProjectId = rs.Primary.Attributes["enterprise_project_id"]
		if authAddr != "" && connAddr != "" {
			if enterpriseProjectId != "" {
				if username != "" && password != "" {
					return fmt.Sprintf("%s/%s/%s/%s/%s/%s", authAddr, connAddr, keyName, username, password, enterpriseProjectId), nil
				}
				return fmt.Sprintf("%s/%s/%s/%s", authAddr, connAddr, keyName, enterpriseProjectId), nil
			}

			if username != "" && password != "" {
				return fmt.Sprintf("%s/%s/%s/%s/%s", authAddr, connAddr, keyName, username, password), nil
			}
			return fmt.Sprintf("%s/%s/%s", authAddr, connAddr, keyName), nil
		}
		return "", fmt.Errorf("missing some attributes: %s/%s/%s", authAddr, connAddr, keyName)
	}
}

func testAccMicroserviceEngineConfiguration_basic_base() string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  id_filter_result = [
    for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == "%[2]s"
  ]
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, acceptance.HW_CSE_MICROSERVICE_ENGINE_ID)
}

func testAccMicroserviceEngineConfiguration_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cse_microservice_engine_configuration" "test" {
  auth_address    = try(local.id_filter_result[0].service_registry_addresses.0.public, null)
  connect_address = try(local.id_filter_result[0].config_center_addresses.0.public, null)
  admin_user      = "root"
  admin_pass      = "%[2]s"

  key        = "%[3]s"
  value_type = "json"
  value      = jsonencode({
    "foo": "bar"
  })

  status                = "enabled"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  tags = {
    owner = "terraform"
  }

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}
`, testAccMicroserviceEngineConfiguration_basic_base(),
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ADMIN_PASSWORD,
		name)
}

func testAccMicroserviceEngineConfiguration_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cse_microservice_engine_configuration" "test" {
  auth_address    = try(local.id_filter_result[0].service_registry_addresses.0.public, null)
  connect_address = try(local.id_filter_result[0].config_center_addresses.0.public, null)
  admin_user      = "root"
  admin_pass      = "%[2]s"

  key        = "%[3]s"
  value_type = "json"
  value      = jsonencode({
    "foo": "baar"
  })

  status                = "disabled"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  tags = {
    owner = "terraform"
  }

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}
`, testAccMicroserviceEngineConfiguration_basic_base(),
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ADMIN_PASSWORD,
		name)
}
