package taurusdb

import (
	"errors"
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

func getTaurusDBSqlControlRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getTaurusDBSqlControlRule: Query the TaurusDB Sql control rule
	var (
		getTaurusDBSqlControlRuleHttpUrl = "v3/{project_id}/instances/{instance_id}/sql-filter/rules"
		getTaurusDBSqlControlRuleProduct = "gaussdb"
	)
	getTaurusDBSqlControlRuleClient, err := cfg.NewServiceClient(getTaurusDBSqlControlRuleProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 4)
	if len(parts) != 4 {
		return nil, errors.New("invalid id format, must be <instance_id>/<node_id>/<sql_type>/<pattern>")
	}
	instanceID := parts[0]
	nodeId := parts[1]
	sqlType := parts[2]
	pattern := parts[3]

	getTaurusDBSqlControlRulePath := getTaurusDBSqlControlRuleClient.Endpoint + getTaurusDBSqlControlRuleHttpUrl
	getTaurusDBSqlControlRulePath = strings.ReplaceAll(getTaurusDBSqlControlRulePath, "{project_id}",
		getTaurusDBSqlControlRuleClient.ProjectID)
	getTaurusDBSqlControlRulePath = strings.ReplaceAll(getTaurusDBSqlControlRulePath, "{instance_id}", instanceID)

	getTaurusDBSqlControlRulePath += buildGetTaurusDBSqlControlRuleQueryParams(nodeId)

	getTaurusDBSqlControlRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getTaurusDBSqlControlRuleResp, err := getTaurusDBSqlControlRuleClient.Request("GET",
		getTaurusDBSqlControlRulePath, &getTaurusDBSqlControlRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving TaurusDB Sql control rule: %s", err)
	}

	getTaurusDBSqlControlRuleRespBody, err := utils.FlattenResponse(getTaurusDBSqlControlRuleResp)
	if err != nil {
		return nil, err
	}

	expression := fmt.Sprintf("sql_filter_rules[?sql_type=='%s']|[0].patterns[?pattern=='%s']|[0].max_concurrency",
		sqlType, pattern)
	maxConcurrency := utils.PathSearch(expression, getTaurusDBSqlControlRuleRespBody, nil)
	if maxConcurrency == nil {
		return nil, errors.New("error get TaurusDB SQL control rule")
	}

	return getTaurusDBSqlControlRuleRespBody, nil
}

func buildGetTaurusDBSqlControlRuleQueryParams(nodeId string) string {
	return fmt.Sprintf("?node_id=%v", nodeId)
}

func TestAccTaurusDBSqlControlRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_taurusdb_sql_control_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTaurusDBSqlControlRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTaurusDBSqlControlRule_basic(name, 20),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_taurusdb_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "node_id",
						"huaweicloud_taurusdb_instance.test", "nodes.0.id"),
					resource.TestCheckResourceAttr(rName, "sql_type", "SELECT"),
					resource.TestCheckResourceAttr(rName, "pattern", "select~from~t1"),
					resource.TestCheckResourceAttr(rName, "max_concurrency", "20"),
				),
			},
			{
				Config: testTaurusDBSqlControlRule_basic(name, 30),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_taurusdb_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "node_id",
						"huaweicloud_taurusdb_instance.test", "nodes.0.id"),
					resource.TestCheckResourceAttr(rName, "sql_type", "SELECT"),
					resource.TestCheckResourceAttr(rName, "pattern", "select~from~t1"),
					resource.TestCheckResourceAttr(rName, "max_concurrency", "30"),
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

func testTaurusDBSqlControlRule_basic(name string, maxConcurrency int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_taurusdb_sql_control_rule" "test" {
  instance_id     = huaweicloud_taurusdb_instance.test.id
  node_id         = huaweicloud_taurusdb_instance.test.nodes[0].id
  sql_type        = "SELECT"
  pattern         = "select~from~t1"
  max_concurrency = %d
}
`, testAccTaurusDBInstanceConfig_basic(name), maxConcurrency)
}
