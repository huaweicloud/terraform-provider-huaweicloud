package gaussdb

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getOpenGaussEipAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getInstanceUrl = "v3/{project_id}/instances?id={instance_id}"
		getEipUrl      = "v3/{project_id}/instances/{instance_id}/public-ips"
		product        = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	getInstancePath := client.Endpoint + getInstanceUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", client.ProjectID)
	getInstancePath = strings.ReplaceAll(getInstancePath, "{instance_id}", state.Primary.Attributes["instance_id"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getInstancePath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	instance := utils.PathSearch("instances[0]", getRespBody, nil)
	if instance == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	node := utils.PathSearch(fmt.Sprintf("nodes[?id == '%s']|[0]", state.Primary.Attributes["node_id"]), instance, nil)
	if node == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	publicIp := utils.PathSearch("public_ip", node, nil)
	if publicIp == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	listEipPath := client.Endpoint + getEipUrl
	listEipPath = strings.ReplaceAll(listEipPath, "{project_id}", client.ProjectID)
	listEipPath = strings.ReplaceAll(listEipPath, "{instance_id}", state.Primary.Attributes["instance_id"])

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listEipPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, err
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return nil, err
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return nil, err
	}

	boundEip := utils.PathSearch(fmt.Sprintf("public_ips[?public_ip_address=='%s']|[0]", publicIp.(string)),
		listRespBody, nil)
	if boundEip == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return listRespBody, nil
}

func TestAccOpenGaussEipAssociate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_opengauss_eip_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOpenGaussEipAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOpenGaussEipAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_gaussdb_opengauss_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "node_id",
						"huaweicloud_gaussdb_opengauss_instance.test", "nodes.0.id"),
					resource.TestCheckResourceAttrPair(rName, "public_ip",
						"huaweicloud_vpc_eip.test", "address"),
					resource.TestCheckResourceAttrPair(rName, "public_ip_id",
						"huaweicloud_vpc_eip.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testOpenGaussEipAssociateImportState(rName),
			},
		},
	})
}

func testOpenGaussEipAssociate_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss_egress" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "egress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = "gaussdb.bs.s6.xlarge.x864.ha"
  name              = "%[2]s"
  password          = "Huangwei!120521"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[3]s"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[2]s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testOpenGaussEipAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_opengauss_eip_associate" "test"{
  instance_id  = huaweicloud_gaussdb_opengauss_instance.test.id
  node_id      = huaweicloud_gaussdb_opengauss_instance.test.nodes.0.id
  public_ip    = huaweicloud_vpc_eip.test.address
  public_ip_id = huaweicloud_vpc_eip.test.id
}
`, testOpenGaussEipAssociate_base(rName))
}

func testOpenGaussEipAssociateImportState(name string) resource.ImportStateIdFunc {
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
