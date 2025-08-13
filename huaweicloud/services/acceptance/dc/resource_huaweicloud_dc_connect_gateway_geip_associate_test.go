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

func getResourceDcDcConnectGatewayGeipAssociateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DC client: %s", err)
	}

	httpUrl := "v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}/binding-global-eips?global_eip_id={global_eip_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connect_gateway_id}", state.Primary.Attributes["connect_gateway_id"])
	getPath = strings.ReplaceAll(getPath, "{global_eip_id}", state.Primary.Attributes["global_eip_id"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DC connect gateway global EIP: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	globalEip := utils.PathSearch("global_eips|[0]", getRespBody, nil)
	if globalEip == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return globalEip, nil
}

func TestAccResourceDcConnectGatewayGeipAssociate_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_dc_connect_gateway_geip_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceDcDcConnectGatewayGeipAssociateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcConnectGatewayId(t)
			acceptance.TestAccPreCheckGlobalEipId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceDcConnectGatewayGeipAssociate_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "connect_gateway_id", acceptance.HW_DC_CONNECT_GATEWAY_ID),
					resource.TestCheckResourceAttr(rName, "global_eip_id", acceptance.HW_GLOBAL_EIP_ID),
					resource.TestCheckResourceAttr(rName, "type", "IP_ADDRESS"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "cidr"),
					resource.TestCheckResourceAttrSet(rName, "address_family"),
					resource.TestCheckResourceAttrSet(rName, "ie_vtep_ip"),
					resource.TestCheckResourceAttrSet(rName, "created_time"),
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

func testResourceDcConnectGatewayGeipAssociate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_connect_gateway_geip_associate" "test" {
  connect_gateway_id = "%[1]s"
  global_eip_id      = "%[2]s"
  type               = "IP_ADDRESS"
}
`, acceptance.HW_DC_CONNECT_GATEWAY_ID, acceptance.HW_GLOBAL_EIP_ID)
}
