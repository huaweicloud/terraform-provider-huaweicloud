package rds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRdsLtsConfigResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		getRdsLtsConfigHttpUrl = "v3/{project_id}/{engine}/instances/logs/lts-configs"
		getRdsLtsConfigProduct = "rds"
	)

	getRdsLtsConfigClient, err := cfg.NewServiceClient(getRdsLtsConfigProduct, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getRdsLtsConfigPath := getRdsLtsConfigClient.Endpoint + getRdsLtsConfigHttpUrl
	getRdsLtsConfigPath = strings.ReplaceAll(getRdsLtsConfigPath, "{project_id}", getRdsLtsConfigClient.ProjectID)
	getRdsLtsConfigPath = strings.ReplaceAll(getRdsLtsConfigPath, "{engine}", state.Primary.Attributes["engine"])
	getRdsLtsConfigPath += fmt.Sprintf("?instance_id=%s", state.Primary.Attributes["instance_id"])

	getRdsLtsConfigResp, err := pagination.ListAllItems(
		getRdsLtsConfigClient,
		"offset",
		getRdsLtsConfigPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS LTS configs: %s", err)
	}

	getRdsLtsConfigRespJson, err := json.Marshal(getRdsLtsConfigResp)
	if err != nil {
		return nil, fmt.Errorf("error marshaling RDS LTS configs: %s", err)
	}

	var getRdsLtsConfigRespBody interface{}
	err = json.Unmarshal(getRdsLtsConfigRespJson, &getRdsLtsConfigRespBody)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling RDS LTS configs: %s", err)
	}

	jsonPath := fmt.Sprintf("instance_lts_configs[0].lts_configs[?log_type=='%s']|[0]", state.Primary.Attributes["log_type"])
	ltsConfig := utils.PathSearch(jsonPath, getRdsLtsConfigRespBody, nil)
	if !utils.PathSearch("enabled", ltsConfig, false).(bool) {
		return nil, fmt.Errorf("error retrieving RDS LTS config: the LTS config can not be found")
	}

	return ltsConfig, nil
}

func TestAccRdsLtsConfig_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_rds_lts_config.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRdsLtsConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsLtsConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "engine", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "log_type", "error_log"),
				),
			},
			{
				Config: testAccRdsLtsConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "engine", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "log_type", "error_log"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRdsLtsConfig_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_rds_lts_config" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  engine        = "mysql"
  log_type      = "error_log"
  lts_group_id  = huaweicloud_lts_group.test.id
  lts_stream_id = huaweicloud_lts_stream.test.id
}`, testAccRdsInstance_mysql_step1(rName), rName)
}

func testAccRdsLtsConfig_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s-update"
  ttl_in_days = 2
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s-update"
}

resource "huaweicloud_rds_lts_config" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  engine        = "mysql"
  log_type      = "error_log"
  lts_group_id  = huaweicloud_lts_group.test.id
  lts_stream_id = huaweicloud_lts_stream.test.id
}`, testAccRdsInstance_mysql_step1(rName), rName)
}
