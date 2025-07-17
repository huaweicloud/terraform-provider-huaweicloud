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

func getDcsInstanceShardBandwidthFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/bandwidths"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	searchExpression := fmt.Sprintf("group_bandwidths[?group_id=='%s']|[0]", state.Primary.Attributes["group_id"])
	groupBandwidth := utils.PathSearch(searchExpression, getRespBody, nil)
	if groupBandwidth == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return groupBandwidth, nil
}

func TestAccDcsInstanceShardBandwidth_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance_shard_bandwidth.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDcsInstanceShardBandwidthFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDcsInstanceShardBandwidth_basic(name, 512),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "group_id",
						"data.huaweicloud_dcs_instance_shards.test", "group_list.0.group_id"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "512"),
					resource.TestCheckResourceAttrSet(rName, "max_bandwidth"),
					resource.TestCheckResourceAttrSet(rName, "assured_bandwidth"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testDcsInstanceShardBandwidth_basic(name, 1024),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bandwidth", "1024"),
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

func testDcsInstanceShardBandwidth_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode       = "ha"
  capacity         = 1
  cpu_architecture = "x86_64"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}

data "huaweicloud_dcs_instance_shards" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
}
`, name)
}

func testDcsInstanceShardBandwidth_basic(name string, bandwidth int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_instance_shard_bandwidth" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
  group_id    = data.huaweicloud_dcs_instance_shards.test.group_list[0].group_id
  bandwidth   = %d
}
`, testDcsInstanceShardBandwidth_base(name), bandwidth)
}
