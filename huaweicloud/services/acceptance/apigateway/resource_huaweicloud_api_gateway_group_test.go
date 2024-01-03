package apigateway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/shared/v1/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApiGatewayGroup_basic(t *testing.T) {
	var resName = "huaweicloud_api_gateway_group.acc_apigw_group"
	rName := fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
	rNameUpdate := rName + "_Update"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApiGatewayGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigwGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayGroupExists(resName),
					resource.TestCheckResourceAttr(
						resName, "name", rName),
					resource.TestCheckResourceAttr(
						resName, "description", "created by acc test"),
				),
			},
			{
				Config: testAccApigwGroup_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayGroupExists(resName),
					resource.TestCheckResourceAttr(
						resName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(
						resName, "description", "updated by acc test"),
				),
			},
		},
	})
}

func testAccCheckApiGatewayGroupDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	apigwClient, err := cfg.ApiGatewayV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating API Gateway client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_api_gateway_group" {
			continue
		}

		_, err := groups.Get(apigwClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("API Gateway Group still exists")
		}
	}

	return nil
}

func testAccCheckApiGatewayGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		apigwClient, err := cfg.ApiGatewayV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating API Gateway client: %s", err)
		}

		found, err := groups.Get(apigwClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("API Gateway Group not found")
		}

		return nil
	}
}

func testAccApigwGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_api_gateway_group" "acc_apigw_group" {
  name        = "%s"
  description = "created by acc test"
}
`, rName)
}

func testAccApigwGroup_update(rNameUpdate string) string {
	return fmt.Sprintf(`
resource "huaweicloud_api_gateway_group" "acc_apigw_group" {
  name        = "%s"
  description = "updated by acc test"
}
`, rNameUpdate)
}
