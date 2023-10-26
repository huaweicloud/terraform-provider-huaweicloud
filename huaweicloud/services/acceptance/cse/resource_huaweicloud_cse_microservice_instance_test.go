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
	token, err := cse.GetAuthorizationToken(state.Primary.Attributes["connect_address"],
		state.Primary.Attributes["admin_user"], state.Primary.Attributes["admin_pass"])
	if err != nil {
		return nil, err
	}

	client := common.NewCustomClient(true, state.Primary.Attributes["connect_address"], "v4", "default")
	return instances.Get(client, state.Primary.Attributes["microservice_id"], state.Primary.ID, token)
}

func TestAccMicroserviceInstance_basic(t *testing.T) {
	var (
		instance     instances.Instance
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cse_microservice_instance.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getMicroserviceInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMicroserviceInstance_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "connect_address",
						"huaweicloud_cse_microservice_engine.test", "service_registry_addresses.0.public"),
					resource.TestCheckResourceAttrPair(resourceName, "microservice_id",
						"huaweicloud_cse_microservice.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "host_name", "localhost"),
					resource.TestCheckResourceAttr(resourceName, "endpoints.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "endpoints.0", "grpc://127.0.1.132:9980"),
					resource.TestCheckResourceAttr(resourceName, "endpoints.1", "rest://127.0.0.111:8081"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttr(resourceName, "properties.nodeIP", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.mode", "push"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.max_retries", "3"),
					resource.TestCheckResourceAttr(resourceName, "data_center.0.name", "dc1"),
					resource.TestCheckResourceAttr(resourceName, "data_center.0.region", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttrPair(resourceName, "data_center.0.availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "status", "UP"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMicroserviceInstanceImportStateIdFunc(),
			},
		},
	})
}

func testAccMicroserviceInstanceImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var connAddr, username, password, microserviceId, instanceId string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_cse_microservice_instance" {
				connAddr = rs.Primary.Attributes["connect_address"]
				microserviceId = rs.Primary.Attributes["microservice_id"]
				username = rs.Primary.Attributes["admin_user"]
				password = rs.Primary.Attributes["admin_pass"]
				instanceId = rs.Primary.ID
			}
		}
		if connAddr != "" && microserviceId != "" && instanceId != "" {
			if username != "" && password != "" {
				return fmt.Sprintf("%s/%s/%s/%s/%s", connAddr, microserviceId, instanceId, username, password), nil
			}
			return fmt.Sprintf("%s/%s/%s", connAddr, microserviceId, instanceId), nil
		}
		return "", fmt.Errorf("resource not found: %s/%s", connAddr, instanceId)
	}
}

func testAccMicroserviceInstance_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cse_microservice" "test" {
  connect_address = huaweicloud_cse_microservice_engine.test.service_registry_addresses.0.public

  name        = "%[2]s"
  app_name    = "%[2]s"
  environment = "development"
  version     = "1.0.1"
  level       = "BACK"

  admin_user = "root"
  admin_pass = huaweicloud_cse_microservice_engine.test.admin_pass
}

resource "huaweicloud_cse_microservice_instance" "test" {
  connect_address = huaweicloud_cse_microservice_engine.test.service_registry_addresses.0.public

  microservice_id = huaweicloud_cse_microservice.test.id
  host_name       = "localhost"
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
    region            = "%[3]s"
    availability_zone = data.huaweicloud_availability_zones.test.names[0]
  }

  admin_user = "root"
  admin_pass = huaweicloud_cse_microservice_engine.test.admin_pass
}
`, testAccMicroserviceEngine_config(rName), rName, acceptance.HW_REGION_NAME)
}
