package cse

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/services"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cse"
)

func getAuthAddress(attributes map[string]string) string {
	if authAddress, ok := attributes["auth_address"]; ok {
		return authAddress
	}
	// Using the connect address as the auth address if its empty.
	// The behavior of the connect address is required.
	return attributes["connect_address"]
}

func getMicroserviceFunc(_ *config.Config, state *terraform.ResourceState) (interface{}, error) {
	token, err := cse.GetAuthorizationToken(getAuthAddress(state.Primary.Attributes),
		state.Primary.Attributes["admin_user"], state.Primary.Attributes["admin_pass"])
	if err != nil {
		return nil, err
	}

	client := common.NewCustomClient(true, state.Primary.Attributes["connect_address"], "v4", "default")
	return services.Get(client, state.Primary.ID, token)
}

// Beforce testing, please bind the EIP and open the access rules according to the resource ducoment appendix.
func TestAccMicroservice_basic(t *testing.T) {
	var (
		service  services.Service
		randName = acceptance.RandomAccResourceName()

		withAuthAddress   = "huaweicloud_cse_microservice.with_auth_address"
		rcWithAuthAddress = acceptance.InitResourceCheck(withAuthAddress, &service, getMicroserviceFunc)

		withoutAuthAddress   = "huaweicloud_cse_microservice.without_auth_address"
		rcWithoutAuthAddress = acceptance.InitResourceCheck(withoutAuthAddress, &service, getMicroserviceFunc)
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
				Config: testAccMicroservice_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					// With auth_address parameter.
					rcWithAuthAddress.CheckResourceExists(),
					resource.TestCheckResourceAttr(withAuthAddress, "name", fmt.Sprintf("%s_with_auth_address", randName)),
					resource.TestCheckResourceAttr(withAuthAddress, "app_name", fmt.Sprintf("%s_with_auth_address", randName)),
					resource.TestCheckResourceAttr(withAuthAddress, "environment", "development"),
					resource.TestCheckResourceAttr(withAuthAddress, "version", "1.0.1"),
					resource.TestCheckResourceAttr(withAuthAddress, "description", "Created by terraform test"),
					resource.TestCheckResourceAttr(withAuthAddress, "level", "BACK"),
					resource.TestCheckResourceAttr(withAuthAddress, "status", "UP"),
					// Without auth_address parameter.
					rcWithoutAuthAddress.CheckResourceExists(),
					resource.TestCheckResourceAttr(withoutAuthAddress, "name", fmt.Sprintf("%s_without_auth_address", randName)),
					resource.TestCheckResourceAttr(withoutAuthAddress, "app_name", fmt.Sprintf("%s_without_auth_address", randName)),
					resource.TestCheckResourceAttr(withoutAuthAddress, "environment", "development"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "version", "1.0.1"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "description", "Created by terraform test"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "level", "BACK"),
					resource.TestCheckResourceAttr(withoutAuthAddress, "status", "UP"),
				),
			},
			{
				ResourceName:      withAuthAddress,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMicroserviceImportStateIdFunc(withAuthAddress),
			},
			{
				ResourceName:      withoutAuthAddress,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMicroserviceImportStateIdFunc(withoutAuthAddress),
				ImportStateVerifyIgnore: []string{
					"auth_address",
				},
			},
		},
	})
}

func testAccMicroserviceImportStateIdFunc(resName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var authAddr, connAddr, addrPart, username, password, microserviceId string
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", resName)
		}

		authAddr = rs.Primary.Attributes["auth_address"]
		connAddr = rs.Primary.Attributes["connect_address"]
		username = rs.Primary.Attributes["admin_user"]
		password = rs.Primary.Attributes["admin_pass"]
		microserviceId = rs.Primary.ID

		addrPart = connAddr
		if authAddr != "" {
			addrPart = fmt.Sprintf("%s/%s", authAddr, addrPart)
		}
		if addrPart != "" && microserviceId != "" {
			if username != "" && password != "" {
				return fmt.Sprintf("%s/%s/%s/%s", addrPart, microserviceId, username, password), nil
			}
			return fmt.Sprintf("%s/%s", addrPart, microserviceId), nil
		}
		return "", fmt.Errorf("missing some attributes: %s/%s", addrPart, microserviceId)
	}
}

func testAccMicroservice_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  id_filter_result = [
    for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == "%[1]s"
  ]
}

resource "huaweicloud_cse_microservice" "with_auth_address" {
  auth_address    = local.id_filter_result[0].service_registry_addresses.0.public
  connect_address = local.id_filter_result[0].service_registry_addresses.0.public

  name        = "%[2]s_with_auth_address"
  app_name    = "%[2]s_with_auth_address"
  environment = "development"
  version     = "1.0.1"
  description = "Created by terraform test"
  level       = "BACK"

  admin_user = "root"
  admin_pass = "%[3]s"

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}

resource "huaweicloud_cse_microservice" "without_auth_address" {
  connect_address = local.id_filter_result[0].service_registry_addresses.0.public

  name        = "%[2]s_without_auth_address"
  app_name    = "%[2]s_without_auth_address"
  environment = "development"
  version     = "1.0.1"
  description = "Created by terraform test"
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
