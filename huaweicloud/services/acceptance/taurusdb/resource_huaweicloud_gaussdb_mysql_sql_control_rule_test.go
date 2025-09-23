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

func getGaussDBSqlControlRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getGaussDBSqlControlRule: Query the GaussDB MySQL Sql control rule
	var (
		getGaussDBSqlControlRuleHttpUrl = "v3/{project_id}/instances/{instance_id}/sql-filter/rules"
		getGaussDBSqlControlRuleProduct = "gaussdb"
	)
	getGaussDBSqlControlRuleClient, err := cfg.NewServiceClient(getGaussDBSqlControlRuleProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<node_id>/<sql_type>/<pattern>")
	}
	instanceID := parts[0]
	nodeId := parts[1]
	sqlType := parts[2]
	pattern := parts[3]

	getGaussDBSqlControlRulePath := getGaussDBSqlControlRuleClient.Endpoint + getGaussDBSqlControlRuleHttpUrl
	getGaussDBSqlControlRulePath = strings.ReplaceAll(getGaussDBSqlControlRulePath, "{project_id}",
		getGaussDBSqlControlRuleClient.ProjectID)
	getGaussDBSqlControlRulePath = strings.ReplaceAll(getGaussDBSqlControlRulePath, "{instance_id}", instanceID)

	getGaussDBSqlControlRulePath += buildGetGaussDBSqlControlRuleQueryParams(nodeId)

	getGaussDBSqlControlRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getGaussDBSqlControlRuleResp, err := getGaussDBSqlControlRuleClient.Request("GET",
		getGaussDBSqlControlRulePath, &getGaussDBSqlControlRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GaussDB MySQL Sql control rule: %s", err)
	}

	getGaussDBSqlControlRuleRespBody, err := utils.FlattenResponse(getGaussDBSqlControlRuleResp)
	if err != nil {
		return nil, err
	}

	expression := fmt.Sprintf("sql_filter_rules[?sql_type=='%s']|[0].patterns[?pattern=='%s']|[0].max_concurrency",
		sqlType, pattern)
	maxConcurrency := utils.PathSearch(expression, getGaussDBSqlControlRuleRespBody, nil)
	if maxConcurrency == nil {
		return nil, fmt.Errorf("error get GaussDB MySQL SQL control rule")
	}

	return getGaussDBSqlControlRuleRespBody, nil
}

func buildGetGaussDBSqlControlRuleQueryParams(nodeId string) string {
	return fmt.Sprintf("?node_id=%v", nodeId)
}

func TestAccGaussDBSqlControlRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_mysql_sql_control_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDBSqlControlRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGaussDBSqlControlRule_basic(name, 20),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_gaussdb_mysql_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "node_id",
						"huaweicloud_gaussdb_mysql_instance.test", "nodes.0.id"),
					resource.TestCheckResourceAttr(rName, "sql_type", "SELECT"),
					resource.TestCheckResourceAttr(rName, "pattern", "select~from~t1"),
					resource.TestCheckResourceAttr(rName, "max_concurrency", "20"),
				),
			},
			{
				Config: testGaussDBSqlControlRule_basic(name, 30),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_gaussdb_mysql_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "node_id",
						"huaweicloud_gaussdb_mysql_instance.test", "nodes.0.id"),
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

func testGaussDBSqlControlRule_basic(name string, maxConcurrency int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_mysql_sql_control_rule" "test" {
  instance_id     = huaweicloud_gaussdb_mysql_instance.test.id
  node_id         = huaweicloud_gaussdb_mysql_instance.test.nodes[0].id
  sql_type        = "SELECT"
  pattern         = "select~from~t1"
  max_concurrency = %d
}
`, testAccGaussDBInstanceConfig_basic(name), maxConcurrency)
}
