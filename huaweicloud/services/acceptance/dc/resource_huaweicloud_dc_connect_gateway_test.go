package dc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceDcDcConnectGatewayFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DC client: %s", err)
	}

	getPath := client.Endpoint + "v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connect_gateway_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DC connect gateway: %s", err)
	}
	return utils.FlattenResponse(getResp)
}

func TestAccResourceDcConnectGateway_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dc_connect_gateway.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceDcDcConnectGatewayFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceDcConnectGateway_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "address_family", "ipv4"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "bgp_asn"),
					resource.TestCheckResourceAttrSet(rName, "current_geip_count"),
					resource.TestCheckResourceAttrSet(rName, "created_time"),
				),
			},
			{
				Config: testResourceDcConnectGateway_basic_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "address_family", "dual"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "bgp_asn"),
					resource.TestCheckResourceAttrSet(rName, "current_geip_count"),
					resource.TestCheckResourceAttrSet(rName, "created_time"),
					resource.TestCheckResourceAttrSet(rName, "updated_time"),
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

func testResourceDcConnectGateway_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_connect_gateway" "test" {
  name           = "%s"
  description    = "test description"
  address_family = "ipv4"
}
`, name)
}

func testResourceDcConnectGateway_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_connect_gateway" "test" {
  name           = "%s"
  description    = ""
  address_family = "dual"
}
`, name)
}
