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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestAccTaurusDBSqlAutoThrottling_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_taurusdb_sql_auto_throttling.test"
	rc := acceptance.InitResourceCheck(rName, &obj, getResourceSqlAutoThrottlingFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBSqlAutoThrottling_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "start_time", "00:00"),
					resource.TestCheckResourceAttr(rName, "end_time", "01:00"),
					resource.TestCheckResourceAttr(rName, "condition", "and"),
					resource.TestCheckResourceAttr(rName, "cpu_usage", "70"),
					resource.TestCheckResourceAttr(rName, "active_sessions", "3"),
					resource.TestCheckResourceAttr(rName, "clear_time", "3"),
					resource.TestCheckResourceAttr(rName, "duration", "2"),
					resource.TestCheckResourceAttr(rName, "max_concurrency", "1000"),
					resource.TestCheckResourceAttr(rName, "retain_sql_rule", "true"),
				),
			},
			{
				Config: testAccTaurusDBSqlAutoThrottling_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "start_time", "02:00"),
					resource.TestCheckResourceAttr(rName, "end_time", "03:00"),
					resource.TestCheckResourceAttr(rName, "condition", "or"),
					resource.TestCheckResourceAttr(rName, "cpu_usage", "80"),
					resource.TestCheckResourceAttr(rName, "active_sessions", "5"),
					resource.TestCheckResourceAttr(rName, "clear_time", "5"),
					resource.TestCheckResourceAttr(rName, "duration", "10"),
					resource.TestCheckResourceAttr(rName, "max_concurrency", "2000"),
					resource.TestCheckResourceAttr(rName, "retain_sql_rule", "false"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"retain_sql_rule"},
			},
		},
	})
}

func getResourceSqlAutoThrottlingFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("gaussdb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	nodeId := state.Primary.Attributes["node_id"]

	return getResourceSqlAutoThrottlingByQuery(client, instanceId, nodeId)
}

func getResourceSqlAutoThrottlingByQuery(client *golangsdk.ServiceClient, instanceId, nodeId string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/auto-sql-limiting"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getOpt.JSONBody = map[string]interface{}{
		"node_ids": []string{nodeId},
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	rules := utils.PathSearch("auto_sql_limiting_rules", getRespBody, make([]interface{}, 0)).([]interface{})
	if len(rules) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	ruleKey := fmt.Sprintf("auto_sql_limiting_rules[?node_id=='%s']", nodeId)
	filteredRules := utils.PathSearch(ruleKey, getRespBody, make([]interface{}, 0)).([]interface{})
	if len(filteredRules) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return filteredRules[0], nil
}

func testAccTaurusDBSqlAutoThrottling_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_sql_auto_throttling" "test" {
  instance_id     = "%[1]s"
  node_id         = "%[2]s"
  start_time      = "00:00"
  end_time        = "01:00"
  condition       = "and"
  cpu_usage       = 70
  active_sessions = 3
  clear_time      = 3
  duration        = 2
  max_concurrency = 1000
  retain_sql_rule = "true"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_NODE_ID)
}

func testAccTaurusDBSqlAutoThrottling_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_sql_auto_throttling" "test" {
  instance_id     = "%[1]s"
  node_id         = "%[2]s"
  start_time      = "02:00"
  end_time        = "03:00"
  condition       = "or"
  cpu_usage       = 80
  active_sessions = 5
  clear_time      = 5
  duration        = 10
  max_concurrency = 2000
  retain_sql_rule = "false"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_NODE_ID)
}
