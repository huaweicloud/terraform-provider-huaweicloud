package cc

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

func getConnectionBandwidthAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/connections"
		product = "cc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CC client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", cfg.DomainID)
	path = strings.ReplaceAll(path, "{central_network_id}", state.Primary.Attributes["central_network_id"])
	path += "?id=" + state.Primary.ID

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", path, &opt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving central network connection: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving central network connection: %s", err)
	}

	connection := utils.PathSearch("central_network_connections", respBody, make([]interface{}, 0))
	if v, ok := connection.([]interface{}); ok && len(v) == 0 {
		return nil, fmt.Errorf("error retrieving central network connection")
	}

	bandwidthType := utils.PathSearch("[0].bandwidth_type", connection, "TestBandwidth").(string)
	if bandwidthType == "TestBandwidth" {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func TestAccCentralNetworkConnectionBandwidthAssociate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_central_network_connection_bandwidth_associate.test"
	baseConfig := testCentralNetworkConnectionBandwidthAssociate_base(name)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getConnectionBandwidthAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCProjectID(t)
			acceptance.TestAccPreCheckCCRegionName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCentralNetworkConnectionBandWidthAssociate_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bandwidth_size", "1"),
					resource.TestCheckResourceAttrPair(rName, "central_network_id",
						"huaweicloud_cc_central_network.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "global_connection_bandwidth_id",
						"huaweicloud_cc_global_connection_bandwidth.b1", "id"),
				),
			},
			{
				Config: testCentralNetworkConnectionBandWidthAssociate_basic_update(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bandwidth_size", "2"),
					resource.TestCheckResourceAttrPair(rName, "central_network_id",
						"huaweicloud_cc_central_network.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "global_connection_bandwidth_id",
						"huaweicloud_cc_global_connection_bandwidth.b2", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCentralNetworkConnectionBandwidthAssociateImportState(rName),
			},
		},
	})
}

func testCentralNetworkConnectionBandwidthAssociateImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["central_network_id"] == "" {
			return "", fmt.Errorf("attribute (central_network_id) of resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["central_network_id"] + "/" + rs.Primary.ID, nil
	}
}

func testCentralNetworkConnectionBandWidthAssociate_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cc_central_network_connections" "test" {
  depends_on = [huaweicloud_cc_central_network_policy_apply.test]

  central_network_id = huaweicloud_cc_central_network.test.id
}

resource "huaweicloud_cc_central_network_connection_bandwidth_associate" test {
  central_network_id             = huaweicloud_cc_central_network.test.id
  connection_id                  = data.huaweicloud_cc_central_network_connections.test.central_network_connections[0].id
  global_connection_bandwidth_id = huaweicloud_cc_global_connection_bandwidth.b1.id
  bandwidth_size                 = 1
}
`, baseConfig)
}

func testCentralNetworkConnectionBandWidthAssociate_basic_update(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cc_central_network_connections" "test" {
  depends_on = [huaweicloud_cc_central_network_policy_apply.test]

  central_network_id = huaweicloud_cc_central_network.test.id
}

resource "huaweicloud_cc_central_network_connection_bandwidth_associate" test {
  central_network_id             = huaweicloud_cc_central_network.test.id
  connection_id                  = data.huaweicloud_cc_central_network_connections.test.central_network_connections[0].id
  global_connection_bandwidth_id = huaweicloud_cc_global_connection_bandwidth.b2.id
  bandwidth_size                 = 2
}
`, baseConfig)
}

func testCentralNetworkConnectionBandwidthAssociate_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "az1" {
  region = "%[1]s"
}

resource "huaweicloud_er_instance" "er1" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.az1.names, 0, 1)

  region                         = "%[1]s"
  name                           = "%[3]s1"
  asn                            = 64512
  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true
}

data "huaweicloud_er_availability_zones" "az2" {
  region = "%[2]s"
}
  
resource "huaweicloud_er_instance" "er2" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.az2.names, 0, 1)
  
  region                         = "%[2]s"
  name                           = "%[3]s2"
  asn                            = 64512
  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true
}

resource "huaweicloud_cc_central_network" "test" {
  name        = "%[3]s"
  description = "This is an accaptance test"
}
 
resource "huaweicloud_cc_central_network_policy" "test" {
  central_network_id = huaweicloud_cc_central_network.test.id
 
  planes {
    associate_er_tables {
      project_id                 = "%[4]s"
      region_id                  = "%[1]s"
      enterprise_router_id       = huaweicloud_er_instance.er1.id
      enterprise_router_table_id = huaweicloud_er_instance.er1.default_association_route_table_id
    }
    associate_er_tables {
      project_id                 = "%[5]s"
      region_id                  = "%[2]s"
      enterprise_router_id       = huaweicloud_er_instance.er2.id
      enterprise_router_table_id = huaweicloud_er_instance.er2.default_association_route_table_id
    }
  }
 
  er_instances {
    project_id           = "%[4]s"
    region_id            = "%[1]s"
    enterprise_router_id = huaweicloud_er_instance.er1.id
  }
  er_instances {
    project_id           = "%[5]s"
    region_id            = "%[2]s"
    enterprise_router_id = huaweicloud_er_instance.er2.id
  }
}

resource "huaweicloud_cc_central_network_policy_apply" "test" {
  central_network_id = huaweicloud_cc_central_network.test.id
  policy_id          = huaweicloud_cc_central_network_policy.test.id
}

resource "huaweicloud_cc_global_connection_bandwidth" "b1" {
  name        = "%[3]s1"
  type        = "Area"  
  bordercross = false
  charge_mode = "bwd"
  size        = 5
  description = "test"
  local_area  = "cn-north-beijing4"
  remote_area = "cn-south-guangzhou"
}
  
resource "huaweicloud_cc_global_connection_bandwidth" "b2" {
  name        = "%[3]s2"
  type        = "Area"  
  bordercross = false
  charge_mode = "bwd"
  size        = 5
  description = "test"
  local_area  = "cn-north-beijing4"
  remote_area = "cn-south-guangzhou"
}	
`, acceptance.HW_REGION_NAME_1, acceptance.HW_REGION_NAME_2,
		name, acceptance.HW_PROJECT_ID_1, acceptance.HW_PROJECT_ID_2)
}
