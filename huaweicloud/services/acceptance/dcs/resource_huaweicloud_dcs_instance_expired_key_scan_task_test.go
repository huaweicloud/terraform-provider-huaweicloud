package dcs

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getInstanceExpiredKeyScanTaskFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/auto-expire/histories"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, err
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return nil, err
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, err
	}

	task := utils.PathSearch(fmt.Sprintf("records[?id=='%s']|[0]", state.Primary.ID), getRespBody, nil)
	if task == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccDcsInstanceExpiredKeyScanTask_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance_expired_key_scan_task.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInstanceExpiredKeyScanTaskFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDcsInstanceExpiredKeyScanTask_baic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_dcs_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "scan_type"),
					resource.TestCheckResourceAttrSet(rName, "num"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "started_at"),
					resource.TestCheckResourceAttrSet(rName, "finished_at"),
				),
			},
		},
	})
}

func testDcsInstanceExpiredKeyScanTask_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 4
  name           = "redis.cluster.xu1.large.r4.4"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 4
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}
`, name)
}

func testDcsInstanceExpiredKeyScanTask_baic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_instance_expired_key_scan_task" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
}
`, testDcsInstanceExpiredKeyScanTask_base(name))
}
