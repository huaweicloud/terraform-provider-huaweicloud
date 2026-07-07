package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/taurusdb"
)

func getResourceProxyEip(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		product    = "gaussdb"
		region     = acceptance.HW_REGION_NAME
		instanceID = state.Primary.Attributes["instance_id"]
		proxyID    = state.Primary.Attributes["proxy_id"]
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB client: %s", err)
	}

	proxyEip, err := taurusdb.GetTaurusDBProxyEip(client, instanceID, proxyID)
	if err != nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return proxyEip, nil
}

func TestAccTaurusDBProxyEipAssociate_basic(t *testing.T) {
	var eip interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_taurusdb_proxy_eip_associate.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getResourceProxyEip,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBProxyEipAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_taurusdb_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "proxy_id",
						"huaweicloud_taurusdb_proxy.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip",
						"huaweicloud_vpc_eip.test", "address"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip_id",
						"huaweicloud_vpc_eip.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testTaurusDBProxyEipAssociateResourceImportState(resourceName),
			},
		},
	})
}

func testAccTaurusDBProxyEipAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_taurusdb_flavors" "test" {
  engine                 = "gaussdb-mysql"
  version                = "8.0"
  availability_zone_mode = "multi"
}

resource "huaweicloud_taurusdb_instance" "test" {
  name                     = "%[2]s"
  password                 = "Test@12345678"
  flavor                   = data.huaweicloud_taurusdb_flavors.test.flavors[0].name
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  enterprise_project_id    = "0"
  master_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  availability_zone_mode   = "multi"
  read_replicas            = 2
}

data "huaweicloud_taurusdb_proxy_flavors" "test" {
  instance_id = huaweicloud_taurusdb_instance.test.id
}

resource "huaweicloud_taurusdb_proxy" "test" {
  instance_id              = huaweicloud_taurusdb_instance.test.id
  flavor                   = data.huaweicloud_taurusdb_proxy_flavors.test.flavor_groups[0].flavors[0].spec_code
  node_num                 = 2
  proxy_name               = "%[2]s"
  proxy_mode               = "readwrite"
  route_mode               = 1
  subnet_id                = huaweicloud_vpc_subnet.test.id
  new_node_auto_add_status = "OFF"
  new_node_weight          = 20

  open_access_control = true
  access_control_type = "white"
  access_control_ip_list {
    ip          = "3.3.3.3"
    description = "test description"
  }
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_taurusdb_proxy_eip_associate" "test" {
  instance_id  = huaweicloud_taurusdb_instance.test.id
  proxy_id     = huaweicloud_taurusdb_proxy.test.id
  public_ip    = huaweicloud_vpc_eip.test.address
  public_ip_id = huaweicloud_vpc_eip.test.id
}
`, common.TestBaseNetwork(rName), rName)
}

func testTaurusDBProxyEipAssociateResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		proxyID := rs.Primary.Attributes["proxy_id"]
		return fmt.Sprintf("%s/%s", instanceID, proxyID), nil
	}
}
