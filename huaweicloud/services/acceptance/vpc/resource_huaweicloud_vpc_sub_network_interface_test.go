package vpc

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

func getSubNetworkInterfaceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v3 client: %s", err)
	}

	getSubNetworkInterfaceHttpUrl := "vpc/sub-network-interfaces/{sub_network_interface_id}"
	getSubNetworkInterfacePath := client.ResourceBaseURL() + getSubNetworkInterfaceHttpUrl
	getSubNetworkInterfacePath = strings.ReplaceAll(getSubNetworkInterfacePath, "{sub_network_interface_id}", state.Primary.ID)

	getSubNetworkInterfaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSubNetworkInterfaceResp, err := client.Request("GET", getSubNetworkInterfacePath, &getSubNetworkInterfaceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving sub network interface: %s", err)
	}

	return utils.FlattenResponse(getSubNetworkInterfaceResp)
}

func TestAccSubNetworkInterface_basic(t *testing.T) {
	var (
		sub_network_interface interface{}
		name                  = acceptance.RandomAccResourceName()
		rName                 = "huaweicloud_vpc_sub_network_interface.test"
		rc                    = acceptance.InitResourceCheck(
			rName,
			&sub_network_interface,
			getSubNetworkInterfaceResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSubNetworkInterface_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "parent_id", "huaweicloud_compute_instance.test", "network.0.port"),
					resource.TestCheckResourceAttr(rName, "description", "created by test acc"),
					resource.TestCheckResourceAttr(rName, "vlan_id", "3002"),
					resource.TestCheckResourceAttr(rName, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "ip_address"),
					resource.TestCheckResourceAttrSet(rName, "mac_address"),
					resource.TestCheckResourceAttrSet(rName, "parent_device_id"),
					resource.TestCheckResourceAttrSet(rName, "vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testSubNetworkInterface_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "parent_id", "huaweicloud_compute_instance.test", "network.0.port"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "vlan_id", "3002"),
					resource.TestCheckResourceAttr(rName, "security_group_ids.#", "1"),
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

func testaccSubNetworkInterface_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_images" "test" {
  architecture = "x86"
  visibility   = "public"
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "computingv3"
  generation        = "c7"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_images.test.images[0].id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, common.TestBaseNetwork(name), name)
}

func testSubNetworkInterface_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_sub_network_interface" "test" {
  subnet_id   = huaweicloud_vpc_subnet.test.id
  parent_id   = huaweicloud_compute_instance.test.network[0].port
  vlan_id     = "3002"
  description = "created by test acc"
}
`, testaccSubNetworkInterface_base(name))
}

func testSubNetworkInterface_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_sub_network_interface" "test" {
  subnet_id   = huaweicloud_vpc_subnet.test.id
  parent_id   = huaweicloud_compute_instance.test.network[0].port
  vlan_id     = "3002"
  description = ""
}
`, testaccSubNetworkInterface_base(name))
}
