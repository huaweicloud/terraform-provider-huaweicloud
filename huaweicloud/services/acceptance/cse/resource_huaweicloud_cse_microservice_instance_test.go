package cse

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cse"
)

func getMicroserviceInstanceFunc(_ *config.Config, state *terraform.ResourceState) (interface{}, error) {
	token, err := cse.GetAuthorizationToken(getAuthAddress(state.Primary.Attributes),
		state.Primary.Attributes["admin_user"], state.Primary.Attributes["admin_pass"])
	if err != nil {
		return nil, err
	}

	client := common.NewCustomClient(true, state.Primary.Attributes["connect_address"], "v4", "default")
	return instances.Get(client, state.Primary.Attributes["microservice_id"], state.Primary.ID, token)
}

// Beforce testing, please bind the EIP and open the access rules according to the resource ducoment appendix.
func TestAccMicroserviceInstance_basic(t *testing.T) {
	var (
		instance instances.Instance
		randName = acceptance.RandomAccResourceName()

		withAuthAddress   = "huaweicloud_cse_microservice_instance.with_auth_address"
		rcWithAuthAddress = acceptance.InitResourceCheck(withAuthAddress, &instance, getMicroserviceInstanceFunc)

		withoutAuthAddress   = "huaweicloud_cse_microservice_instance.without_auth_address"
		rcWithoutAuthAddress = acceptance.InitResourceCheck(withoutAuthAddress, &instance, getMicroserviceInstanceFunc)
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
				Config: testAccMicroserviceInstance_basic(randName),
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
					resource.TestCheckResourceAttr(withAuthAddress, "status", "UP"),
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
					resource.TestCheckResourceAttr(withoutAuthAddress, "status", "UP"),
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
		var authAddr, connAddr, addrPart, microserviceId, username, password, instanceId string
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", resName)
		}

		authAddr = rs.Primary.Attributes["auth_address"]
		connAddr = rs.Primary.Attributes["connect_address"]
		microserviceId = rs.Primary.Attributes["microservice_id"]
		username = rs.Primary.Attributes["admin_user"]
		password = rs.Primary.Attributes["admin_pass"]
		instanceId = rs.Primary.ID

		addrPart = connAddr
		if authAddr != "" {
			addrPart = fmt.Sprintf("%s/%s", authAddr, addrPart)
		}
		if addrPart != "" && microserviceId != "" && instanceId != "" {
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
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  id_filter_result = [
    for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == "%[1]s"
  ]
}

resource "huaweicloud_cse_microservice" "test" {
  auth_address    = local.id_filter_result[0].service_registry_addresses[0].public
  connect_address = local.id_filter_result[0].service_registry_addresses[0].public

  name        = "%[2]s"
  app_name    = "%[2]s"
  environment = "development"
  version     = "1.0.1"
  level       = "BACK"

  admin_user = "root"
  admin_pass = "%[3]s"

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}
`, acceptance.HW_CSE_MICROSERVICE_ENGINE_ID,
		name,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ADMIN_PASSWORD)
}

func testAccMicroserviceInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cse_microservice_instance" "with_auth_address" {
  auth_address    = local.id_filter_result[0].service_registry_addresses[0].public
  connect_address = local.id_filter_result[0].service_registry_addresses[0].public

  microservice_id = huaweicloud_cse_microservice.test.id
  host_name       = "localhost_with_auth_address"
  endpoints       = ["grpc://127.0.1.132:9980", "rest://127.0.0.111:8081"]
  version         = "1.0.1"

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
    region            = "%[2]s"
    availability_zone = data.huaweicloud_availability_zones.test.names[0]
  }

  admin_user = "root"
  admin_pass = "%[3]s"

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}

resource "huaweicloud_cse_microservice_instance" "without_auth_address" {
  connect_address = local.id_filter_result[0].service_registry_addresses[0].public

  microservice_id = huaweicloud_cse_microservice.test.id
  host_name       = "localhost_without_auth_address"
  endpoints       = ["grpc://127.0.1.132:9980", "rest://127.0.0.111:8081"]
  version         = "1.0.1"

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
    region            = "%[2]s"
    availability_zone = data.huaweicloud_availability_zones.test.names[0]
  }

  admin_user = "root"
  admin_pass = "%[3]s"

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}
`, testAccMicroserviceInstance_base(name),
		acceptance.HW_REGION_NAME,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ADMIN_PASSWORD)
}
