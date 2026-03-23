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

func getMicroserviceInstanceFunc(_ *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		// Querying microservice instances requires building a client based on the microservice engine's connection address,
		// which does not use IAM authentication.
		client   = common.NewCustomClient(true, state.Primary.Attributes["connect_address"])
		authInfo = cse.MicroserviceEngineAuthInfo{
			AuthAddress:         getAuthAddress(state.Primary.Attributes),
			AdminUser:           state.Primary.Attributes["admin_user"],
			AdminPass:           state.Primary.Attributes["admin_pass"],
			EnterpriseProjectId: state.Primary.Attributes["enterprise_project_id"],
		}
		microserviceId = state.Primary.Attributes["microservice_id"]
		instanceId     = state.Primary.ID
	)

	return cse.GetInstance(client, authInfo, microserviceId, instanceId)
}

// Beforce testing, please bind the EIP and open the access rules according to the resource ducoment appendix.
func TestAccMicroserviceInstance_basic(t *testing.T) {
	var (
		obj interface{}

		withAuthAddress   = "huaweicloud_cse_microservice_instance.with_auth_address"
		rcWithAuthAddress = acceptance.InitResourceCheck(withAuthAddress, &obj, getMicroserviceInstanceFunc)

		withoutAuthAddress   = "huaweicloud_cse_microservice_instance.without_auth_address"
		rcWithoutAuthAddress = acceptance.InitResourceCheck(withoutAuthAddress, &obj, getMicroserviceInstanceFunc)

		randName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSEMicroserviceEngineID(t)
			acceptance.TestAccPreCheckCSEMicroserviceEngineAdminPassword(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithAuthAddress.CheckResourceDestroy(),
			rcWithoutAuthAddress.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccMicroserviceInstance_basic_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					// With auth_address parameter.
					rcWithAuthAddress.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(withAuthAddress, "microservice_id",
						"huaweicloud_cse_microservice.test", "id"),
					resource.TestCheckResourceAttr(withAuthAddress, "host_name", "localhost_with_auth_address"),
					resource.TestCheckResourceAttr(withAuthAddress, "endpoints.#", "2"),
					resource.TestCheckResourceAttr(withAuthAddress, "endpoints.0", "grpc://127.0.1.132:9980"),
					resource.TestCheckResourceAttr(withAuthAddress, "endpoints.1", "rest://127.0.0.111:8081"),
					resource.TestCheckResourceAttr(withAuthAddress, "version", "1.0.1"),
					resource.TestCheckResourceAttr(withAuthAddress, "properties.nodeIP", "127.0.0.1"),
					resource.TestCheckResourceAttr(withAuthAddress, "health_check.0.mode", "push"),
					resource.TestCheckResourceAttr(withAuthAddress, "health_check.0.interval", "30"),
					resource.TestCheckResourceAttr(withAuthAddress, "health_check.0.max_retries", "3"),
					resource.TestCheckResourceAttr(withAuthAddress, "data_center.0.name", "dc1"),
					resource.TestCheckResourceAttr(withAuthAddress, "data_center.0.region", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttrPair(withAuthAddress, "data_center.0.availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(withAuthAddress, "status", "DOWN"),
					// Without auth_address parameter.
					rcWithoutAuthAddress.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(withAuthAddress, "microservice_id",
						"huaweicloud_cse_microservice.test", "id"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "host_name", "localhost_without_auth_address"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "endpoints.#", "2"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "endpoints.0", "grpc://127.0.1.132:9980"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "endpoints.1", "rest://127.0.0.111:8081"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "version", "1.0.1"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "properties.nodeIP", "127.0.0.1"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "health_check.0.mode", "push"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "health_check.0.interval", "30"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "health_check.0.max_retries", "3"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "data_center.0.name", "dc1"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "data_center.0.region", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttrPair(withoutAuthAddress, "data_center.0.availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "status", "DOWN"),
				),
			},
			{
				Config: testAccMicroserviceInstance_basic_step2(randName),
				Check: resource.ComposeTestCheckFunc(
					rcWithAuthAddress.CheckResourceExists(),
					resource.TestCheckResourceAttr(withAuthAddress, "status", "UP"),
					rcWithoutAuthAddress.CheckResourceExists(),
					resource.TestCheckResourceAttr(withoutAuthAddress, "status", "UP"),
				),
			},
			{
				Config: testAccMicroserviceInstance_basic_step3(randName),
				Check: resource.ComposeTestCheckFunc(
					rcWithAuthAddress.CheckResourceExists(),
					resource.TestCheckResourceAttr(withAuthAddress, "status", "OUTOFSERVICE"),
					rcWithoutAuthAddress.CheckResourceExists(),
					resource.TestCheckResourceAttr(withoutAuthAddress, "status", "OUTOFSERVICE"),
				),
			},
			{
				ResourceName:      withAuthAddress,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMicroserviceInstanceImportStateIdFunc(withAuthAddress),
			},
			{
				ResourceName:      withoutAuthAddress,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMicroserviceInstanceImportStateIdFunc(withoutAuthAddress),
				ImportStateVerifyIgnore: []string{
					"auth_address",
				},
			},
		},
	})
}

func testAccMicroserviceInstanceImportStateIdFunc(resName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var authAddr, connAddr, addrPart, microserviceId, username, password, enterpriseProjectId, instanceId string
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", resName)
		}

		authAddr = rs.Primary.Attributes["auth_address"]
		connAddr = rs.Primary.Attributes["connect_address"]
		microserviceId = rs.Primary.Attributes["microservice_id"]
		username = rs.Primary.Attributes["admin_user"]
		password = rs.Primary.Attributes["admin_pass"]
		enterpriseProjectId = rs.Primary.Attributes["enterprise_project_id"]
		instanceId = rs.Primary.ID

		addrPart = connAddr
		if authAddr != "" {
			addrPart = fmt.Sprintf("%s/%s", authAddr, addrPart)
		}
		if addrPart != "" && microserviceId != "" && instanceId != "" {
			if enterpriseProjectId != "" {
				if username != "" && password != "" {
					return fmt.Sprintf("%s/%s/%s/%s/%s/%s", addrPart, microserviceId, instanceId, username, password, enterpriseProjectId), nil
				}
				return fmt.Sprintf("%s/%s/%s/%s", addrPart, microserviceId, instanceId, enterpriseProjectId), nil
			}

			if username != "" && password != "" {
				return fmt.Sprintf("%s/%s/%s/%s/%s", addrPart, microserviceId, instanceId, username, password), nil
			}
			return fmt.Sprintf("%s/%s/%s", addrPart, microserviceId, instanceId), nil
		}
		return "", fmt.Errorf("missing some attributes: %s/%s/%s", addrPart, microserviceId, instanceId)
	}
}

func testAccMicroserviceInstance_base(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  id_filter_result = [
    for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == "%[2]s"
  ]
}

resource "huaweicloud_cse_microservice" "test" {
  auth_address    = try(local.id_filter_result[0].service_registry_addresses[0].public, null)
  connect_address = try(local.id_filter_result[0].service_registry_addresses[0].public, null)
  admin_user      = "root"
  admin_pass      = "%[3]s"

  name                  = "%[4]s"
  app_name              = "%[4]s"
  environment           = "development"
  version               = "1.0.1"
  level                 = "BACK"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ID,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ADMIN_PASSWORD,
		name)
}

func testAccMicroserviceInstance_basic_with_status(name, status string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cse_microservice_instance" "with_auth_address" {
  auth_address    = try(local.id_filter_result[0].service_registry_addresses[0].public, null)
  connect_address = try(local.id_filter_result[0].service_registry_addresses[0].public, null)
  admin_user      = "root"
  admin_pass      = "%[2]s"

  microservice_id       = huaweicloud_cse_microservice.test.id
  host_name             = "localhost_with_auth_address"
  endpoints             = ["grpc://127.0.1.132:9980", "rest://127.0.0.111:8081"]
  version               = "1.0.1"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  properties = {
    "nodeIP" = "127.0.0.1"
  }

  health_check {
    mode        = "push"
    interval    = 30
    max_retries = 3
  }

  data_center {
    name              = "dc1"
    region            = "%[3]s"
    availability_zone = data.huaweicloud_availability_zones.test.names[0]
  }

  status = "%[4]s"

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}

resource "huaweicloud_cse_microservice_instance" "without_auth_address" {
  auth_address    = try(local.id_filter_result[0].service_registry_addresses[0].public, null)
  connect_address = try(local.id_filter_result[0].service_registry_addresses[0].public, null)
  admin_user      = "root"
  admin_pass      = "%[2]s"

  microservice_id       = huaweicloud_cse_microservice.test.id
  host_name             = "localhost_without_auth_address"
  endpoints             = ["grpc://127.0.1.132:9980", "rest://127.0.0.111:8081"]
  version               = "1.0.1"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  properties = {
    "nodeIP" = "127.0.0.1"
  }

  health_check {
    mode        = "push"
    interval    = 30
    max_retries = 3
  }

  data_center {
    name              = "dc1"
    region            = "%[3]s"
    availability_zone = data.huaweicloud_availability_zones.test.names[0]
  }

  status = "%[4]s"

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}
`, testAccMicroserviceInstance_base(name),
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ADMIN_PASSWORD,
		acceptance.HW_REGION_NAME,
		status)
}

func testAccMicroserviceInstance_basic_step1(name string) string {
	return testAccMicroserviceInstance_basic_with_status(name, "DOWN")
}

func testAccMicroserviceInstance_basic_step2(name string) string {
	return testAccMicroserviceInstance_basic_with_status(name, "UP")
}

func testAccMicroserviceInstance_basic_step3(name string) string {
	return testAccMicroserviceInstance_basic_with_status(name, "OUTOFSERVICE")
}
