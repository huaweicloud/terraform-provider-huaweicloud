package vpn

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCustomerGatewayResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getCustomerGateway: Query the VPN customer gateway detail
	var (
		getCustomerGatewayHttpUrl = "v5/{project_id}/customer-gateways/{id}"
		getCustomerGatewayProduct = "vpn"
	)
	getCustomerGatewayClient, err := config.NewServiceClient(getCustomerGatewayProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CustomerGateway Client: %s", err)
	}

	getCustomerGatewayPath := getCustomerGatewayClient.Endpoint + getCustomerGatewayHttpUrl
	getCustomerGatewayPath = strings.ReplaceAll(getCustomerGatewayPath, "{project_id}", getCustomerGatewayClient.ProjectID)
	getCustomerGatewayPath = strings.ReplaceAll(getCustomerGatewayPath, "{id}", state.Primary.ID)

	getCustomerGatewayOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getCustomerGatewayResp, err := getCustomerGatewayClient.Request("GET", getCustomerGatewayPath, &getCustomerGatewayOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CustomerGateway: %s", err)
	}
	return utils.FlattenResponse(getCustomerGatewayResp)
}

func TestAccCustomerGateway_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	nameUpdate := name + "-update"
	rName := "huaweicloud_vpn_customer_gateway.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCustomerGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCustomerGateway_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ip", "192.168.1.1"),
				),
			},
			{
				Config: testCustomerGateway_basic(nameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", nameUpdate),
					resource.TestCheckResourceAttr(rName, "ip", "192.168.1.1"),
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

func testCustomerGateway_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_customer_gateway" "test" {
  name = "%s"
  ip   = "172.16.1.1"
}
`, name)
}
