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

func getVirtualInterfaceSwitchoverResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/dcaas/switchover-test?resource_id={resource_id}"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DC Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{resource_id}", state.Primary.Attributes["resource_id"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	searchPath := fmt.Sprintf("switchover_test_records[?id=='%s']|[0]", state.Primary.ID)
	switchover := utils.PathSearch(searchPath, getRespBody, nil)
	if switchover == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return switchover, nil
}

func TestAccVirtualInterfaceSwitchover_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dc_virtual_interface_switchover.test"
	updateResourceName := "huaweicloud_dc_virtual_interface_switchover.update"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getVirtualInterfaceSwitchoverResourceFunc,
	)

	urc := acceptance.InitResourceCheck(
		updateResourceName,
		&obj,
		getVirtualInterfaceSwitchoverResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcDirectConnection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testVirtualInterfaceSwitchover_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "resource_id",
						"huaweicloud_dc_virtual_interface.test", "id"),
					resource.TestCheckResourceAttr(rName, "operation", "shutdown"),
					resource.TestCheckResourceAttr(rName, "resource_type", "virtual_interface"),
					resource.TestCheckResourceAttrSet(rName, "start_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttrSet(rName, "operate_status"),
				),
			},
			{
				Config: testVirtualInterfaceSwitchover_update(name),
				Check: resource.ComposeTestCheckFunc(
					urc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(updateResourceName, "resource_id",
						"huaweicloud_dc_virtual_interface.test", "id"),
					resource.TestCheckResourceAttr(updateResourceName, "operation", "undo_shutdown"),
					resource.TestCheckResourceAttr(updateResourceName, "resource_type", "virtual_interface"),
					resource.TestCheckResourceAttrSet(updateResourceName, "start_time"),
					resource.TestCheckResourceAttrSet(updateResourceName, "end_time"),
					resource.TestCheckResourceAttrSet(updateResourceName, "operate_status"),
				),
			},
		},
	})
}

func testVirtualInterfaceSwitchover_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_dc_virtual_gateway" "test" {
  vpc_id = huaweicloud_vpc.test.id
  name   = "%[1]s"

  local_ep_group = [
    huaweicloud_vpc.test.cidr,
  ]
}

resource "huaweicloud_dc_virtual_interface" "test" {
  direct_connect_id = "%[2]s"
  vgw_id            = huaweicloud_dc_virtual_gateway.test.id
  name              = "%[1]s"
  description       = "Created by acc test"
  type              = "private"
  route_mode        = "static"
  vlan              = 80
  bandwidth         = 5
  enable_bfd        = true
  enable_nqa        = false

  remote_ep_group = [
    "1.1.1.0/30",
  ]

  address_family       = "ipv4"
  local_gateway_v4_ip  = "1.1.1.1/30"
  remote_gateway_v4_ip = "1.1.1.2/30"
}
`, name, acceptance.HW_DC_DIRECT_CONNECT_ID)
}

func testVirtualInterfaceSwitchover_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dc_virtual_interface_switchover" "test" {
  resource_id   = huaweicloud_dc_virtual_interface.test.id
  operation     = "shutdown"
  resource_type = "virtual_interface"
}
`, testVirtualInterfaceSwitchover_base(name))
}

func testVirtualInterfaceSwitchover_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dc_virtual_interface_switchover" "update" {
  resource_id   = huaweicloud_dc_virtual_interface.test.id
  operation     = "undo_shutdown"
  resource_type = "virtual_interface"
}
`, testVirtualInterfaceSwitchover_base(name))
}
