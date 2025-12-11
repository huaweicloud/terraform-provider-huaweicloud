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

func getResourceDcVifPeerDetectionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DC client: %s", err)
	}

	getPath := client.Endpoint + "v3/{project_id}/dcaas/vif-peer-detections/{id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DC vif peer detection: %s", err)
	}
	return utils.FlattenResponse(getResp)
}

func TestAccResourceDcVifPeerDetection_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dc_vif_peer_detection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceDcVifPeerDetectionFunc,
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
				Config: testResourceDcVifPeerDetection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "vif_peer_id",
						"data.huaweicloud_dc_virtual_interfaces.test",
						"virtual_interfaces.0.vif_peers.0.id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "start_time"),
					resource.TestCheckResourceAttrSet(rName, "loss_ratio"),
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

func testResourceDcVifPeerDetection_base(name string) string {
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

data "huaweicloud_dc_virtual_interfaces" "test" {
  virtual_interface_id = huaweicloud_dc_virtual_interface.test.id
}
`, name, acceptance.HW_DC_DIRECT_CONNECT_ID)
}

func testResourceDcVifPeerDetection_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dc_vif_peer_detection" "test" {
  vif_peer_id = data.huaweicloud_dc_virtual_interfaces.test.virtual_interfaces[0].vif_peers[0].id
}
`, testResourceDcVifPeerDetection_base(name))
}
