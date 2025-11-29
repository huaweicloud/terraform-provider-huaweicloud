package esw

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

func getEswConnection(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/l2cg/instances/{instance_id}/connections/{connection_id}"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ESW client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getPath = strings.ReplaceAll(getPath, "{connection_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func TestAccEswConnection_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_esw_connection.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getEswConnection,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEswConnection_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_esw_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "virsubnet_id",
						"huaweicloud_vpc_subnet.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ips.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "remote_infos.0.segmentation_id", "9999"),
					resource.TestCheckResourceAttrPair(resourceName, "remote_infos.0.tunnel_ip",
						"huaweicloud_esw_instance.test", "tunnel_info.0.tunnel_ip"),
					resource.TestCheckResourceAttrPair(resourceName, "remote_infos.0.tunnel_port",
						"huaweicloud_esw_instance.test", "tunnel_info.0.tunnel_port"),

					resource.TestCheckResourceAttrSet(resourceName, "remote_infos.0.tunnel_type"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccEswConnection_basic_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_esw_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "virsubnet_id",
						"huaweicloud_vpc_subnet.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ips.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "remote_infos.0.segmentation_id", "9999"),
					resource.TestCheckResourceAttrPair(resourceName, "remote_infos.0.tunnel_ip",
						"huaweicloud_esw_instance.test", "tunnel_info.0.tunnel_ip"),
					resource.TestCheckResourceAttrPair(resourceName, "remote_infos.0.tunnel_port",
						"huaweicloud_esw_instance.test", "tunnel_info.0.tunnel_port"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testEswConnectionResourceImportState(resourceName),
				ImportStateVerifyIgnore: []string{"status", "updated_at"},
			},
		},
	})
}

func testAccEswConnection_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_esw_flavors" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  count = 2

  name       = "%[1]s-${count.index}"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.${count.index}.0/24"
  gateway_ip = "192.168.${count.index}.1"
}

resource "huaweicloud_esw_instance" "test" {
  name        = "%[1]s"
  flavor_ref  = data.huaweicloud_esw_flavors.test.flavors[0].name
  ha_mode     = "ha"
  description = "test description"

  availability_zones {
    primary = data.huaweicloud_esw_flavors.test.flavors.0.available_zones[0]
    standby = data.huaweicloud_esw_flavors.test.flavors.0.available_zones[1]
  }

  tunnel_info {
    vpc_id       = huaweicloud_vpc.test.id
    virsubnet_id = huaweicloud_vpc_subnet.test[0].id
    tunnel_ip    = "192.168.0.192"
  }

  charge_infos {
    charge_mode = "postPaid"
  }
}
`, rName)
}

func testAccEswConnection_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_esw_connection" "test" {
  instance_id  = huaweicloud_esw_instance.test.id
  name         = "%[2]s"
  vpc_id       = huaweicloud_vpc.test.id
  virsubnet_id = huaweicloud_vpc_subnet.test[1].id
  fixed_ips    = ["192.168.1.80", "192.168.1.100"]

  remote_infos {
    segmentation_id = 9999
    tunnel_ip       = huaweicloud_esw_instance.test.tunnel_info[0].tunnel_ip
    tunnel_port     = huaweicloud_esw_instance.test.tunnel_info[0].tunnel_port
  }
}
`, testAccEswConnection_base(rName), rName)
}

func testAccEswConnection_basic_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_esw_connection" "test" {
  instance_id  = huaweicloud_esw_instance.test.id
  name         = "%[2]s"
  vpc_id       = huaweicloud_vpc.test.id
  virsubnet_id = huaweicloud_vpc_subnet.test[1].id
  fixed_ips    = ["192.168.1.80", "192.168.1.100"]

  remote_infos {
    segmentation_id = 9999
    tunnel_ip       = huaweicloud_esw_instance.test.tunnel_info[0].tunnel_ip
    tunnel_port     = huaweicloud_esw_instance.test.tunnel_info[0].tunnel_port
  }
}
`, testAccEswConnection_base(rName), updateName)
}

func testEswConnectionResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceID, rs.Primary.ID), nil
	}
}
