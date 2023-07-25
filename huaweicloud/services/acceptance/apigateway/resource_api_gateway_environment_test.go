package apigateway

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/shared/v1/environments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getEnvironmentFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApiGatewayV1Client(acceptance.HW_REGION_NAME)
	envId := state.Primary.ID
	log.Printf("[DEBUG] env id is : %s", envId)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client %s", err)
	}

	envs, err := environments.List(client, environments.ListOpts{
		EnvName: state.Primary.Attributes["name"],
	})
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] List of shared APIG environments: %#v", envs)
	for i, v := range envs {
		if v.Id == envId {
			return &envs[i], nil
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccEnvironment_basic(t *testing.T) {
	var env environments.Environment
	rName := "huaweicloud_api_gateway_environment.test_env"
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(rName, &env, getEnvironmentFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironment_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by acc test"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testAccEnvironment_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "updated by acc test"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEnvironment_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_api_gateway_environment" "test_env" {
  name        = "%s"
  description = "created by acc test"
}
`, rName)
}

func testAccEnvironment_update(rNameUpdate string) string {
	return fmt.Sprintf(`
resource "huaweicloud_api_gateway_environment" "test_env" {
  name        = "%s"
  description = "updated by acc test"
}
`, rNameUpdate)
}
