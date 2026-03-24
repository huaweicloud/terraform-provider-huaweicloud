package ddm

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

func getInstanceGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/groups"
		product = "ddm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDM client: %s", err)
	}

	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", state.Primary.Attributes["instance_id"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	offset := 0
	var group interface{}
	for {
		getPath := getBasePath + buildPageQueryParams(offset)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, err
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}
		group = utils.PathSearch(fmt.Sprintf("group_list|[?id=='%s']|[0]", state.Primary.ID), getRespBody, nil)
		if group != nil {
			break
		}
		groups := utils.PathSearch("group_list", getRespBody, make([]interface{}, 0)).([]interface{})
		offset += len(groups)
		totalCount := utils.PathSearch("total_count", getRespBody, float64(0)).(float64)
		if offset >= int(totalCount) {
			break
		}
	}
	if group == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return group, nil
}

func buildPageQueryParams(offset int) string {
	return fmt.Sprintf("?limit=100&offset=%d", offset)
}

func TestAccInstanceGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ddm_instance_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInstanceGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInstanceGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "rw"),
					resource.TestCheckResourceAttrSet(rName, "endpoint"),
					resource.TestCheckResourceAttrSet(rName, "is_load_balance"),
					resource.TestCheckResourceAttrSet(rName, "is_default_group"),
					resource.TestCheckResourceAttrSet(rName, "cpu_num_per_node"),
					resource.TestCheckResourceAttrSet(rName, "mem_num_per_node"),
					resource.TestCheckResourceAttrSet(rName, "architecture"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testInstanceGroupImportState(rName),
				ImportStateVerifyIgnore: []string{"flavor_id", "nodes"},
			},
		},
	})
}

func testInstanceGroup_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.huaweicloud_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.huaweicloud_ddm_engines.test.engines[0].id
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  admin_user        = "test_user_1"
  admin_password    = "test_password_123"

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  lifecycle {
    ignore_changes = [
      node_num
    ]
  }
}
`, testDdmInstance_base(name), name)
}

func testInstanceGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ddm_instance_group" "test" {
  instance_id = huaweicloud_ddm_instance.test.id
  name        = "%[2]s"
  type        = "rw"
  flavor_id   = data.huaweicloud_ddm_flavors.test.flavors[0].id

  nodes {
    available_zone = "cn-north-4a"
    subnet_id      = huaweicloud_vpc_subnet.test.id
  }
  nodes {
    available_zone = "cn-north-4a"
    subnet_id      = huaweicloud_vpc_subnet.test.id
  }

  lifecycle {
    ignore_changes = [
      nodes
    ]
  }
}
`, testInstanceGroup_base(name), name)
}

func testInstanceGroupImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("the resource (%s) not found: %s", name, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceId, rs.Primary.ID), nil
	}
}
