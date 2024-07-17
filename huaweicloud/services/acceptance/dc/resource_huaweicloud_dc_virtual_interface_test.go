package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dc/v3/interfaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getVirtualInterfaceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DcV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DC v3 client: %s", err)
	}

	return interfaces.Get(client, state.Primary.ID)
}

func TestAccVirtualInterface_basic(t *testing.T) {
	var (
		vif interfaces.VirtualInterface

		rName      = "huaweicloud_dc_virtual_interface.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		vlan       = acctest.RandIntRange(1, 3999)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&vif,
		getVirtualInterfaceFunc,
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
				Config: testAccVirtualInterface_basic(name, vlan),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "direct_connect_id", acceptance.HW_DC_DIRECT_CONNECT_ID),
					resource.TestCheckResourceAttrPair(rName, "vgw_id", "huaweicloud_dc_virtual_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "type", "private"),
					resource.TestCheckResourceAttr(rName, "route_mode", "static"),
					resource.TestCheckResourceAttr(rName, "vlan", fmt.Sprintf("%v", vlan)),
					resource.TestCheckResourceAttr(rName, "bandwidth", "5"),
					resource.TestCheckResourceAttr(rName, "enable_bfd", "true"),
					resource.TestCheckResourceAttr(rName, "enable_nqa", "false"),
					resource.TestCheckResourceAttr(rName, "remote_ep_group.0", "1.1.1.0/30"),
					resource.TestCheckResourceAttr(rName, "address_family", "ipv4"),
					resource.TestCheckResourceAttr(rName, "local_gateway_v4_ip", "1.1.1.1/30"),
					resource.TestCheckResourceAttr(rName, "remote_gateway_v4_ip", "1.1.1.2/30"),
					resource.TestCheckResourceAttrSet(rName, "device_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),

					// Computed Attribute `vif_peers`
					resource.TestCheckResourceAttr(rName, "vif_peers.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.address_family"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.bgp_asn"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.bgp_route_limit"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.device_id"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.enable_bfd"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.enable_nqa"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.id"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.local_gateway_ip"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.name"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.receive_route_num"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.remote_ep_group.#"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.remote_gateway_ip"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.route_mode"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.status"),
					resource.TestCheckResourceAttrSet(rName, "vif_peers.0.vif_id"),
				),
			},
			{
				Config: testAccVirtualInterface_update1(updateName, vlan),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "direct_connect_id", acceptance.HW_DC_DIRECT_CONNECT_ID),
					resource.TestCheckResourceAttrPair(rName, "vgw_id", "huaweicloud_dc_virtual_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "type", "private"),
					resource.TestCheckResourceAttr(rName, "route_mode", "static"),
					resource.TestCheckResourceAttr(rName, "vlan", fmt.Sprintf("%v", vlan)),
					resource.TestCheckResourceAttr(rName, "bandwidth", "10"),
					resource.TestCheckResourceAttr(rName, "enable_bfd", "false"),
					resource.TestCheckResourceAttr(rName, "enable_nqa", "true"),
					resource.TestCheckResourceAttr(rName, "remote_ep_group.0", "1.1.1.0/30"),
					resource.TestCheckResourceAttr(rName, "remote_ep_group.1", "1.1.2.0/30"),
					resource.TestCheckResourceAttr(rName, "address_family", "ipv4"),
					resource.TestCheckResourceAttr(rName, "local_gateway_v4_ip", "1.1.1.1/30"),
					resource.TestCheckResourceAttr(rName, "remote_gateway_v4_ip", "1.1.1.2/30"),
					resource.TestCheckResourceAttrSet(rName, "device_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttr(rName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_update"),
				),
			},
			{
				Config: testAccVirtualInterface_update2(updateName, vlan),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "direct_connect_id", acceptance.HW_DC_DIRECT_CONNECT_ID),
					resource.TestCheckResourceAttr(rName, "enable_bfd", "true"),
					resource.TestCheckResourceAttr(rName, "enable_nqa", "false"),
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

func TestAccVirtualInterface_acrossTenant(t *testing.T) {
	var (
		vif interfaces.VirtualInterface

		rName = "huaweicloud_dc_virtual_interface.test"
		name  = acceptance.RandomAccResourceName()
		vlan  = acctest.RandIntRange(1, 3999)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&vif,
		getVirtualInterfaceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcDirectConnection(t)
			acceptance.TestAccPreCheckTargetTenantDcVGW(t)
			acceptance.TestAccPreCheckDcResourceTenant(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualInterface_acrossTenant(name, vlan),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "direct_connect_id", acceptance.HW_DC_DIRECT_CONNECT_ID),
					resource.TestCheckResourceAttr(rName, "vgw_id", acceptance.HW_DC_TARGET_TENANT_VGW_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "type", "private"),
					resource.TestCheckResourceAttr(rName, "route_mode", "static"),
					resource.TestCheckResourceAttr(rName, "vlan", fmt.Sprintf("%v", vlan)),
					resource.TestCheckResourceAttr(rName, "bandwidth", "10"),
					resource.TestCheckResourceAttr(rName, "enable_bfd", "false"),
					resource.TestCheckResourceAttr(rName, "enable_nqa", "true"),
					resource.TestCheckResourceAttr(rName, "resource_tenant_id", acceptance.HW_DC_RESOURCE_TENANT_ID),
					resource.TestCheckResourceAttr(rName, "remote_ep_group.0", "1.1.1.0/30"),
					resource.TestCheckResourceAttr(rName, "remote_ep_group.1", "1.1.2.0/30"),
					resource.TestCheckResourceAttr(rName, "address_family", "ipv4"),
					resource.TestCheckResourceAttr(rName, "local_gateway_v4_ip", "1.1.1.1/30"),
					resource.TestCheckResourceAttr(rName, "remote_gateway_v4_ip", "1.1.1.2/30"),
					resource.TestCheckResourceAttrSet(rName, "device_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_tenant_id"},
			},
		},
	})
}

func testAccVirtualInterface_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_dc_virtual_gateway" "test" {
  vpc_id      = huaweicloud_vpc.test.id
  name        = "%[1]s"
  description = "Created by acc test"

  local_ep_group = [
    huaweicloud_vpc.test.cidr,
  ]
}
`, name)
}

func testAccVirtualInterface_basic(name string, vlan int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dc_virtual_interface" "test" {
  direct_connect_id = "%[2]s"
  vgw_id            = huaweicloud_dc_virtual_gateway.test.id
  name              = "%[3]s"
  description       = "Created by acc test"
  type              = "private"
  route_mode        = "static"
  vlan              = %[4]d
  bandwidth         = 5
  enable_bfd        = true
  enable_nqa        = false

  remote_ep_group = [
    "1.1.1.0/30",
  ]

  address_family       = "ipv4"
  local_gateway_v4_ip  = "1.1.1.1/30"
  remote_gateway_v4_ip = "1.1.1.2/30"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccVirtualInterface_base(name), acceptance.HW_DC_DIRECT_CONNECT_ID, name, vlan)
}

func testAccVirtualInterface_update1(name string, vlan int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dc_virtual_interface" "test" {
  direct_connect_id = "%[2]s"
  vgw_id            = huaweicloud_dc_virtual_gateway.test.id
  name              = "%[3]s"
  type              = "private"
  route_mode        = "static"
  vlan              = %[4]d
  bandwidth         = 10
  enable_bfd        = false
  enable_nqa        = true

  remote_ep_group = [
    "1.1.1.0/30",
    "1.1.2.0/30",
  ]

  address_family       = "ipv4"
  local_gateway_v4_ip  = "1.1.1.1/30"
  remote_gateway_v4_ip = "1.1.1.2/30"

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, testAccVirtualInterface_base(name), acceptance.HW_DC_DIRECT_CONNECT_ID, name, vlan)
}

func testAccVirtualInterface_update2(name string, vlan int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dc_virtual_interface" "test" {
  direct_connect_id = "%[2]s"
  vgw_id            = huaweicloud_dc_virtual_gateway.test.id
  name              = "%[3]s"
  type              = "private"
  route_mode        = "static"
  vlan              = %[4]d
  bandwidth         = 10
  enable_bfd        = true
  enable_nqa        = false

  remote_ep_group = [
    "1.1.1.0/30",
    "1.1.2.0/30",
  ]

  address_family       = "ipv4"
  local_gateway_v4_ip  = "1.1.1.1/30"
  remote_gateway_v4_ip = "1.1.1.2/30"

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, testAccVirtualInterface_base(name), acceptance.HW_DC_DIRECT_CONNECT_ID, name, vlan)
}

func testAccVirtualInterface_acrossTenant(name string, vlan int) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_virtual_interface" "test" {
  direct_connect_id = "%[1]s"
  vgw_id            = "%[2]s"
  name              = "%[3]s"
  type              = "private"
  route_mode        = "static"
  vlan              = %[4]d
  bandwidth         = 10
  enable_bfd        = false
  enable_nqa        = true

  resource_tenant_id = "%[5]s"

  remote_ep_group = [
    "1.1.1.0/30",
    "1.1.2.0/30",
  ]

  address_family       = "ipv4"
  local_gateway_v4_ip  = "1.1.1.1/30"
  remote_gateway_v4_ip = "1.1.1.2/30"
}
`, acceptance.HW_DC_DIRECT_CONNECT_ID, acceptance.HW_DC_TARGET_TENANT_VGW_ID, name, vlan,
		acceptance.HW_DC_RESOURCE_TENANT_ID)
}
