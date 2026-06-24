package geminidb

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

func getGeminidbDBInstanceLtsLogAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/logs/lts-configs"
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	logType := state.Primary.Attributes["log_type"]

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGaussdbInstanceLtsLogAssociateBodyParams(instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	listPath := fmt.Sprintf("instance_lts_configs[?instance.id=='%s']|[0].lts_configs", instanceId)
	ltsList := utils.PathSearch(listPath, getRespBody, []interface{}{}).([]interface{})

	filterPath := fmt.Sprintf("[?log_type=='%s']|[0]", logType)
	matchedConfig := utils.PathSearch(filterPath, ltsList, nil)
	if matchedConfig == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	enabled := utils.PathSearch("enabled", matchedConfig, false).(bool)
	if !enabled {
		return nil, golangsdk.ErrDefault404{}
	}
	return matchedConfig, nil
}

func buildGaussdbInstanceLtsLogAssociateBodyParams(instanceId string) string {
	res := ""
	res = fmt.Sprintf("%s&instance_id=%v", res, instanceId)
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func TestAccResourceGeminidbInstanceLtsLogAssociate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	nameUpdate := acceptance.RandomAccResourceName()
	rName := "huaweicloud_geminidb_instance_lts_log_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGeminidbDBInstanceLtsLogAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGeminidbInstanceLtsLogAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttr(rName, "log_type", "slow_log"),
					resource.TestCheckResourceAttrSet(rName, "lts_group_id"),
					resource.TestCheckResourceAttrSet(rName, "lts_stream_id"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
				),
			},
			{
				Config: testAccResourceGeminidbInstanceLtsLogAssociate_update(nameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "log_type", "slow_log"),
					resource.TestCheckResourceAttrPair(rName, "lts_stream_id", "huaweicloud_lts_stream.test_new", "id"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
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

func testAccResourceGeminidbInstanceLtsLogAssociate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_geminidb_instance_lts_log_associate" "test" {
  instance_id   = "%[1]s"
  log_type      = "slow_log"
  lts_group_id  = huaweicloud_lts_group.test.id
  lts_stream_id = huaweicloud_lts_stream.test.id
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID, name)
}

func testAccResourceGeminidbInstanceLtsLogAssociate_update(newName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test_new" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_geminidb_instance_lts_log_associate" "test" {
  instance_id   = "%[1]s"
  log_type      = "slow_log"
  lts_group_id  = huaweicloud_lts_group.test.id
  lts_stream_id = huaweicloud_lts_stream.test_new.id
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID, newName)
}
