package dcs

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

func getDcsNodePriorityConfigFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/logical-nodes"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	nodeId := state.Primary.Attributes["node_id"]

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	node := utils.PathSearch(fmt.Sprintf("nodes[?node_id=='%s']|[0]", nodeId), respBody, nil)
	if node == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return node, nil
}

func TestAccDcsNodePriorityConfig_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_dcs_node_priority_config.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDcsNodePriorityConfigFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsNodePriorityConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "slave_priority_weight", "50"),
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttrSet(rName, "group_id"),
					resource.TestCheckResourceAttrSet(rName, "node_id"),
					resource.TestCheckResourceAttrSet(rName, "logical_node_id"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "az_code"),
					resource.TestCheckResourceAttrSet(rName, "node_role"),
					resource.TestCheckResourceAttrSet(rName, "node_type"),
					resource.TestCheckResourceAttrSet(rName, "node_ip"),
					resource.TestCheckResourceAttrSet(rName, "priority_weight"),
				),
			},
			{
				Config: testAccDcsNodePriorityConfig_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "slave_priority_weight", "80"),
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

func testAccDcsNodePriorityConfig_basic(name string) string {
	return fmt.Sprintf(`
%s

locals {
  replication_list = data.huaweicloud_dcs_instance_shards.test.group_list[0].replication_list
}

resource "huaweicloud_dcs_node_priority_config" "test" {
  instance_id           = huaweicloud_dcs_instance.test.id
  group_id              = data.huaweicloud_dcs_instance_shards.test.group_list[0].group_id
  node_id               = [for v in local.replication_list : v.node_id if v.replication_role == "slave"][0]
  slave_priority_weight = 50
}
`, testDcsInstanceNodeIpRemove_base(name))
}

func testAccDcsNodePriorityConfig_update(name string) string {
	return fmt.Sprintf(`
%s

locals {
  replication_list = data.huaweicloud_dcs_instance_shards.test.group_list[0].replication_list
}

resource "huaweicloud_dcs_node_priority_config" "test" {
  instance_id           = huaweicloud_dcs_instance.test.id
  group_id              = data.huaweicloud_dcs_instance_shards.test.group_list[0].group_id
  node_id               = [for v in local.replication_list : v.node_id if v.replication_role == "slave"][0]
  slave_priority_weight = 80
}
`, testDcsInstanceNodeIpRemove_base(name))
}
