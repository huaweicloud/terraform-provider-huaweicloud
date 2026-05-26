package taurusdb

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

func getResourceInstanceNodeConfig(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/details?instance_ids={instance_id}"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])

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
	node := utils.PathSearch(fmt.Sprintf("instances[0].nodes[?id=='%s']|[0]", state.Primary.ID), getRespBody, nil)
	if node == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return node, nil
}

func TestAccTaurusDBInstanceNodeConfig_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	nodeName := acceptance.RandomAccResourceNameWithDash()
	updateNodeName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_taurusdb_instance_node_config.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceInstanceNodeConfig,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBInstanceNodeConfig_basic(rName, nodeName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", nodeName),
					resource.TestCheckResourceAttr(resourceName, "priority", "3"),
				),
			},
			{
				Config: testAccTaurusDBInstanceNodeConfig_basic_update(rName, updateNodeName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateNodeName),
					resource.TestCheckResourceAttr(resourceName, "priority", "5"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testTaurusDBInstanceNodeConfigImportState(resourceName),
			},
		},
	})
}

func testAccTaurusDBInstanceNodeConfig_base(rName string) string {
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
`, common.TestBaseNetwork(rName), rName)
}

func testAccTaurusDBInstanceNodeConfig_basic(rName, nodeName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_instance_node_config" "test" {
  instance_id = huaweicloud_taurusdb_instance.test.id
  node_id     = huaweicloud_taurusdb_instance.test.nodes[0].id
  name        = "%[2]s"
  priority    = "3"
}`, testAccTaurusDBInstanceNodeConfig_base(rName), nodeName)
}

func testAccTaurusDBInstanceNodeConfig_basic_update(rName, updateNodeName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_instance_node_config" "test" {
  instance_id = huaweicloud_taurusdb_instance.test.id
  node_id     = huaweicloud_taurusdb_instance.test.nodes[0].id
  name        = "%[2]s"
  priority    = "5"
}`, testAccTaurusDBInstanceNodeConfig_base(rName), updateNodeName)
}

func testTaurusDBInstanceNodeConfigImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("the resource (%s) not found: %s", name, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		nodeId := rs.Primary.Attributes["node_id"]
		return fmt.Sprintf("%s/%s", instanceId, nodeId), nil
	}
}
