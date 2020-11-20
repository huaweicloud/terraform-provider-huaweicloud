package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/apigw/groups"
)

func TestAccApiGatewayGroup_basic(t *testing.T) {
	var resName = "huaweicloud_api_gateway_group.acc_apigw_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApiGatewayGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigwGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayGroupExists(resName),
					resource.TestCheckResourceAttr(
						resName, "name", "acc_apigw_group_1"),
					resource.TestCheckResourceAttr(
						resName, "description", "created by acc test"),
				),
			},
			{
				Config: testAccApigwGroup_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayGroupExists(resName),
					resource.TestCheckResourceAttr(
						resName, "name", "acc_apigw_group_update"),
					resource.TestCheckResourceAttr(
						resName, "description", "updated by acc test"),
				),
			},
		},
	})
}

func testAccCheckApiGatewayGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	apigwClient, err := config.apiGatewayV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud api gateway client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_api_gateway_group" {
			continue
		}

		_, err := groups.Get(apigwClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("api gateway group still exists")
		}
	}

	return nil
}

func testAccCheckApiGatewayGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		apigwClient, err := config.apiGatewayV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud api gateway client: %s", err)
		}

		found, err := groups.Get(apigwClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("apigateway group not found")
		}

		return nil
	}
}

const testAccApigwGroup_basic = `
resource "huaweicloud_api_gateway_group" "acc_apigw_group" {
	name = "acc_apigw_group_1"
	description = "created by acc test"
}
`
const testAccApigwGroup_update = `
resource "huaweicloud_api_gateway_group" "acc_apigw_group" {
	name = "acc_apigw_group_update"
	description = "updated by acc test"
}
`
