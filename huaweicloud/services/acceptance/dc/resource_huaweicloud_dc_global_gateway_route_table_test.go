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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceGlobalGatewayRouteTableFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DC client: %s", err)
	}

	getPath := client.Endpoint + "v3/{project_id}/dcaas/gdgw/{gdgw_id}/routetables"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{gdgw_id}", state.Primary.Attributes["gdgw_id"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DC connect gateway route table: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	routeTable := utils.PathSearch(fmt.Sprintf("gdgw_routetables[?id=='%s']|[0]", state.Primary.ID), getRespBody, nil)
	if routeTable == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccResourceDcGlobalGatewayRouteTable_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dc_global_gateway_route_table.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceGlobalGatewayRouteTableFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcDirectConnection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceDcGlobalGatewayRouteTable_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "gdgw_id",
						"huaweicloud_dc_global_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "type", "vif_peer"),
					resource.TestCheckResourceAttr(rName, "destination", "2.2.3.0/30"),
					resource.TestCheckResourceAttrPair(rName, "nexthop",
						"huaweicloud_dc_virtual_interface.test", "vif_peers.0.id"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttrSet(rName, "obtain_mode"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "address_family"),
				),
			},
			{
				Config: testResourceDcGlobalGatewayRouteTable_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "gdgw_id",
						"huaweicloud_dc_global_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "type", "vif_peer"),
					resource.TestCheckResourceAttr(rName, "destination", "2.2.3.0/30"),
					resource.TestCheckResourceAttrPair(rName, "nexthop",
						"huaweicloud_dc_virtual_interface.test", "vif_peers.0.id"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrSet(rName, "obtain_mode"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "address_family"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGlobalGatewayRouteTableImportStateFunc(rName),
			},
		},
	})
}

func testResourceDcGlobalGatewayRouteTable_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dc_virtual_gateway" "test" {
  vpc_id = huaweicloud_vpc.test.id
  name   = "%[2]s"

  local_ep_group = [
    huaweicloud_vpc.test.cidr,
  ]
}

resource "huaweicloud_dc_global_gateway" "test" {
  name           = "%[2]s"
  description    = "test description"
  bgp_asn        = 10
  address_family = "ipv4"
}

resource "huaweicloud_dc_virtual_interface" "test" {
  direct_connect_id = "%[3]s"
  vgw_id            = huaweicloud_dc_virtual_gateway.test.id
  name              = "%[2]s"
  type              = "private"
  route_mode        = "static"
  vlan              = 70
  bandwidth         = 10
  asn               = 200
  enable_bfd        = true
  enable_nqa        = false
  service_type      = "GDGW"
  gateway_id        = huaweicloud_dc_global_gateway.test.id

  remote_ep_group = [
    "1.1.1.0/30",
  ]

  address_family       = "ipv4"
  local_gateway_v4_ip  = "1.1.1.1/30"
  remote_gateway_v4_ip = "1.1.1.2/30"

  lifecycle {
    ignore_changes = [
      remote_ep_group,
    ]
  }
}
`, common.TestVpc(name), name, acceptance.HW_DC_DIRECT_CONNECT_ID)
}

func testResourceDcGlobalGatewayRouteTable_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dc_global_gateway_route_table" "test" {
  gdgw_id     = huaweicloud_dc_global_gateway.test.id
  type        = "vif_peer"
  destination = "2.2.3.0/30"
  nexthop     = huaweicloud_dc_virtual_interface.test.vif_peers[0].id
  description = "test description"
}
`, testResourceDcGlobalGatewayRouteTable_base(name))
}

func testResourceDcGlobalGatewayRouteTable_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dc_global_gateway_route_table" "test" {
  gdgw_id     = huaweicloud_dc_global_gateway.test.id
  type        = "vif_peer"
  destination = "2.2.3.0/30"
  nexthop     = huaweicloud_dc_virtual_interface.test.vif_peers[0].id
  description = ""
}
`, testResourceDcGlobalGatewayRouteTable_base(name))
}

func testAccGlobalGatewayRouteTableImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["gdgw_id"] == "" {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["gdgw_id"], rs.Primary.ID), nil
	}
}
